package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	SUCCESS                      = "SUCCESS"
	INVALID_CONTRACT             = "INVALID_CONTRACT"
	PARENT_DISTRIBUTOR_NOT_FOUND = "PARENT_DISTRIBUTOR_NOT_FOUND"
	DISTRIBUTOR_NOT_FOUND        = "DISTRIBUTOR_NOT_FOUND"
)

func TestApplyContractSelfValidation(t *testing.T) {
	ts := SetupIntegrationTest(t)
	defer CleanupTest(t, ts)

	tests := []struct {
		name                 string
		contract             string
		expectedStatusCode   int
		expectedStatus       bool
		expectedResponseCode string
		// expectedError  string
	}{

		{
			name: "Empty distributor",
			contract: `Permissions for 
INCLUDE: PT
INCLUDE: US`,
			expectedStatusCode:   http.StatusBadRequest,
			expectedStatus:       false,
			expectedResponseCode: INVALID_CONTRACT,
		},
		{
			name: "Invalid heading format",
			contract: `Permissifsfdsfdsf DISTRIBUTOR2 < DISTRIBUTOR1
INCLUDE: PT
INCLUDE: US`,
			expectedStatusCode:   http.StatusBadRequest,
			expectedStatus:       false,
			expectedResponseCode: INVALID_CONTRACT,
		},
		{
			name: "Invalid line in contract",
			contract: `Permissions for DISTRIBUTOR1
INCLUDE: PT
INCLUDE: US
BLA BLA`,
			expectedStatusCode:   http.StatusBadRequest,
			expectedStatus:       false,
			expectedResponseCode: INVALID_CONTRACT,
		},
		{
			name: "Non-existent parent distributor",
			contract: `Permissions for DISTRIBUTOR2 < DISTRIBUTOR121323
INCLUDE: IN
INCLUDE: US`,
			expectedStatusCode:   http.StatusNotFound,
			expectedStatus:       false,
			expectedResponseCode: PARENT_DISTRIBUTOR_NOT_FOUND,
		},
		{
			name: "Invalid region",
			contract: `Permissions for DISTRIBUTOR1
INCLUDE: IN
INCLUDE: US
EXCLUDE: KAA-IN`,
			expectedStatusCode:   http.StatusNotFound,
			expectedStatus:       false,
			expectedResponseCode: "REGION_NOT_FOUND",
		},
		{
			name: "Duplicate distributor",
			contract: `Permissions for DISTRIBUTOR1 < DISTRIBUTOR1
INCLUDE: IN
INCLUDE: US`,
			expectedStatusCode:   http.StatusBadRequest,
			expectedStatus:       false,
			expectedResponseCode: INVALID_CONTRACT,
		},
		{
			name: "Valid contract with only excludes",
			contract: `Permissions for DISTRIBUTOR4
EXCLUDE: KA-IN
EXCLUDE: CENAI-TN-IN`,
			expectedStatusCode:   http.StatusBadRequest,
			expectedStatus:       false,
			expectedResponseCode: INVALID_CONTRACT,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		})
	}
}

func TestApplyContractWithParentExistence(t *testing.T) {
	ts := SetupIntegrationTest(t)
	defer CleanupTest(t, ts)

	// First create parent distributor
	parentContract := `Permissions for DISTRIBUTOR1
INCLUDE: IN
INCLUDE: US`

	req := httptest.NewRequest("POST", "/permission/contract", strings.NewReader(parentContract))
	req.Header.Set("Content-Type", "text/plain")
	resp, err := ts.App.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	tests := []struct {
		name                 string
		contract             string
		expectedStatusCode   int
		expectedStatus       bool
		expectedResponseCode string
	}{
		{

			name: "Valid contract with existing parent",
			contract: `Permissions for DISTRIBUTOR2 < DISTRIBUTOR1
INCLUDE: IN
INCLUDE: US
EXCLUDE: KA-IN
EXCLUDE: CENAI-TN-IN`,
			expectedStatusCode:   http.StatusOK,
			expectedStatus:       true,
			expectedResponseCode: SUCCESS,
		},
		{

			name: "Valid contract with existing parent",
			contract: `Permissions for DISTRIBUTOR2 < DISTRIBUTOR1
INCLUDE: IN
INCLUDE: US
EXCLUDE: KA-IN
EXCLUDE: CENAI-TN-IN`,
			expectedStatusCode:   http.StatusOK,
			expectedStatus:       true,
			expectedResponseCode: SUCCESS,
		},
		{

			name: "Valid contract with non-existing parent",
			contract: `Permissions for DISTRIBUTOR2 < DISTRIBUTOR243434
INCLUDE: IN
INCLUDE: US
EXCLUDE: KA-IN
EXCLUDE: CENAI-TN-IN`,
			expectedStatusCode:   http.StatusNotFound,
			expectedStatus:       false,
			expectedResponseCode: PARENT_DISTRIBUTOR_NOT_FOUND,
		},
		{
			name: "Valid contract without parent",
			contract: `Permissions for DISTRIBUTOR3
INCLUDE: IN
INCLUDE: US`,
			expectedStatusCode:   http.StatusOK,
			expectedStatus:       true,
			expectedResponseCode: SUCCESS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		})
	}
}
