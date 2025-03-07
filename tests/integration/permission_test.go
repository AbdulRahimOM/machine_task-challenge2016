package test

import (
	"challenge16/internal/dto"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDistributorPermissionsFlow(t *testing.T) {
	ts := SetupIntegrationTest(t)
	defer CleanupTest(t, ts)

	tests := []struct {
		name                 string
		contract             string
		expectedStatusCode   int
		expectedStatus       bool
		expectedResponseCode string

		recipientDistributor      string
		expectedStatusCodeInGet   int
		expectedStatusInGet       bool
		expectedResponseCodeInGet string
		expectedIncluded          []string
		expectedExcluded          []string
	}{
		{
			name: "Initial contract",
			contract: `Permissions for DISTRIBUTOR1
INCLUDE: IN
INCLUDE: US
EXCLUDE: KA-IN
EXCLUDE: CENAI-TN-IN`,
			expectedStatusCode:   http.StatusOK,
			expectedStatus:       true,
			expectedResponseCode: SUCCESS,

			recipientDistributor:      "DISTRIBUTOR1",
			expectedStatusCodeInGet:   http.StatusOK,
			expectedStatusInGet:       true,
			expectedResponseCodeInGet: SUCCESS,
			expectedIncluded:          []string{"IN", "US"},
			expectedExcluded:          []string{"KA-IN", "CENAI-TN-IN"},
		},
		{
			name: "Sub contract with same inclusion and exclusion",
			contract: `Permissions for DISTRIBUTOR2 < DISTRIBUTOR1
INCLUDE: IN
INCLUDE: US
EXCLUDE: KA-IN
EXCLUDE: CENAI-TN-IN`,
			expectedStatusCode:   http.StatusOK,
			expectedStatus:       true,
			expectedResponseCode: SUCCESS,

			recipientDistributor:      "DISTRIBUTOR2",
			expectedStatusCodeInGet:   http.StatusOK,
			expectedStatusInGet:       true,
			expectedResponseCodeInGet: SUCCESS,
			expectedIncluded:          []string{"IN", "US"},
			expectedExcluded:          []string{"KA-IN", "CENAI-TN-IN"},
		},
		{
			name: "Sub contract with syntax mistake",
			contract: `Permissionss for DISTRIBUTOR3 < DISTRIBUTOR1
INCLUDE: IN
INCLUDE: US
EXCLUDE: KA-IN
EXCLUDE: CENAI-TN-IN`,
			expectedStatusCode:   http.StatusBadRequest,
			expectedStatus:       false,
			expectedResponseCode: INVALID_CONTRACT,

			recipientDistributor:      "DISTRIBUTOR3",
			expectedStatusCodeInGet:   http.StatusNotFound,
			expectedStatusInGet:       false,
			expectedResponseCodeInGet: DISTRIBUTOR_NOT_FOUND,
		},
		{
			name: "Sub contract with different inclusion and exclusion",
			contract: `Permissions for DISTRIBUTOR3 < DISTRIBUTOR1
INCLUDE: IN
EXCLUDE: KA-IN`,
			expectedStatusCode:   http.StatusOK,
			expectedStatus:       true,
			expectedResponseCode: SUCCESS,

			recipientDistributor:      "DISTRIBUTOR3",
			expectedStatusCodeInGet:   http.StatusOK,
			expectedStatusInGet:       true,
			expectedResponseCodeInGet: SUCCESS,
			expectedIncluded:          []string{"IN"},
			expectedExcluded:          []string{"KA-IN", "CENAI-TN-IN"},
		},
		{
			name: "Sub contract with extra inclusion and exclusion",
			contract: `Permissions for DISTRIBUTOR4 < DISTRIBUTOR1
INCLUDE: IN
INCLUDE: US
INCLUDE: PA
EXCLUDE: KA-IN
EXCLUDE: CENAI-TN-IN
EXCLUDE: GJ-IN`,
			expectedStatusCode:   http.StatusOK,
			expectedStatus:       true,
			expectedResponseCode: SUCCESS,

			recipientDistributor:      "DISTRIBUTOR4",
			expectedStatusCodeInGet:   http.StatusOK,
			expectedStatusInGet:       true,
			expectedResponseCodeInGet: SUCCESS,
			expectedIncluded:          []string{"IN", "US"},
			expectedExcluded:          []string{"KA-IN", "CENAI-TN-IN", "GJ-IN"},
		},
		{
			name: "Misc fresh contract",
			contract: `Permissions for DISTRIBUTOR5
INCLUDE: IN
INCLUDE: PA
EXCLUDE: KA-IN`,
			expectedStatusCode:   http.StatusOK,
			expectedStatus:       true,
			expectedResponseCode: SUCCESS,

			recipientDistributor:      "DISTRIBUTOR5",
			expectedStatusCodeInGet:   http.StatusOK,
			expectedStatusInGet:       true,
			expectedResponseCodeInGet: SUCCESS,
			expectedIncluded:          []string{"IN", "PA"},
			expectedExcluded:          []string{"KA-IN"},
		},
		{
			name: "Giving sub-contract to distributor having some permissions",
			contract: `Permissions for DISTRIBUTOR5 < DISTRIBUTOR1
INCLUDE: IN
EXCLUDE: AP-IN`,
			expectedStatusCode:   http.StatusOK,
			expectedStatus:       true,
			expectedResponseCode: SUCCESS,

			recipientDistributor:      "DISTRIBUTOR5",
			expectedStatusCodeInGet:   http.StatusOK,
			expectedStatusInGet:       true,
			expectedResponseCodeInGet: SUCCESS,
			expectedIncluded:          []string{"IN", "PA"},
			expectedExcluded:          []string{"KA-IN"}, //CENA-TN-IN is not excluded because it is already permitted for DISTRIBUTOR5, so exclusion in contract is ignored
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// POST request to create a contract
			req := httptest.NewRequest("POST", "/permission/contract", strings.NewReader(tt.contract))
			req.Header.Set("Content-Type", "text/plain")

			resp, err := ts.App.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			var response Response
			err = json.Unmarshal(body, &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, response.Status)
			assert.Equal(t, tt.expectedResponseCode, response.ResponseCode)

			if !tt.expectedStatus == response.Status {
				t.Skip("Skipping GET request test: No valid distributor found in contract")
			}

			// GET request to verify stored permissions
			getReq := httptest.NewRequest("GET", fmt.Sprintf("/permission/%s?type=json", tt.recipientDistributor), nil)
			getResp, err := ts.App.Test(getReq)
			assert.NoError(t, err)

			body, err = io.ReadAll(getResp.Body)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatusCodeInGet, getResp.StatusCode)

			var getResponse struct {
				Status   bool        `json:"status"`
				RespCode string      `json:"resp_code"`
				Data     interface{} `json:"data,omitempty"`
				Error    string      `json:"error,omitempty"`
			}
			err = json.Unmarshal(body, &getResponse)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedResponseCodeInGet, getResponse.RespCode)
			assert.Equal(t, tt.expectedStatusInGet, getResponse.Status)

			if tt.expectedStatusInGet == getResponse.Status {
				dataBytes, err := json.Marshal(getResponse.Data)
				assert.NoError(t, err)

				var permissionData dto.GetPermissionsData
				err = json.Unmarshal(dataBytes, &permissionData)
				assert.NoError(t, err)

				// Compare included and excluded regions
				assert.ElementsMatch(t, tt.expectedIncluded, permissionData.Included)
				assert.ElementsMatch(t, tt.expectedExcluded, permissionData.Excluded)
			}

		})
	}
}
