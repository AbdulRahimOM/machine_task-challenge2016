package handler

import (
	"challenge16/internal/dto"
	"challenge16/internal/response"
	"challenge16/validation"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type checkPermissionRequest struct {
	Distributor  string `query:"distributor" validate:"required"`
	RegionString string `query:"region" validate:"required"`
}

func (h *handler) CheckIfDistributionIsAllowed(c *fiber.Ctx) error {
	req := new(checkPermissionRequest)

	if ok, err := validation.BindAndValidateURLQueryRequest(c, req); !ok {
		return err
	}

	resp := h.databank.CheckIfDistributionIsAllowed(req.Distributor, req.RegionString)
	return resp.WriteToJSON(c)
}

func (h *handler) AllowDistribution(c *fiber.Ctx) error {
	req := new(struct {
		RegionString string `json:"region" validate:"required"`
		Distributor  string `json:"distributor" validate:"required"`
	})

	if ok, err := validation.BindAndValidateJSONRequest(c, req); !ok {
		return err
	}

	resp := h.databank.MarkInclusion(req.Distributor, req.RegionString)
	return resp.WriteToJSON(c)
}

func (h *handler) DisallowDistribution(c *fiber.Ctx) error {
	req := new(struct {
		RegionString string `json:"region" validate:"required"`
		Distributor  string `json:"distributor" validate:"required"`
	})

	if ok, err := validation.BindAndValidateJSONRequest(c, req); !ok {
		return err
	}

	resp := h.databank.MarkExclusion(req.Distributor, req.RegionString)
	return resp.WriteToJSON(c)
}

func (h *handler) ApplyContract(c *fiber.Ctx) error {
	contractText := string(c.Body())
	contract, err := getContractData(contractText)
	if err != nil {
		return response.Response{
			HttpStatusCode: 400,
			ResponseCode:   "INVALID_CONTRACT",
			Error:          err,
		}.WriteToJSON(c)
	}
	resp := h.databank.ApplyContract(*contract)
	return resp.WriteToJSON(c)
}

func getContractData(contractText string) (*dto.Contract, error) {
	//Example contract:
	/*
		Permissions for DISTRIBUTOR1
		INCLUDE: IN
		INCLUDE: UN
		EXCLUDE: KA-IN
		EXCLUDE: CENAI-TN-IN
	*/

	//or

	/*
		Permissions for DISTRIBUTOR1 < DISTRIBUTOR2 < DISTRIBUTOR3
		INCLUDE: YADGR-KA-IN
	*/
	var (
		contract = dto.Contract{
			Permissions: dto.Permissions{
				IncludedCountries: make(map[string]bool),
				IncludedProvinces: make(map[string]map[string]bool),
				IncludedCities:    make(map[string]map[string]map[string]bool),
				ExcludedProvinces: make(map[string]map[string]bool),
				ExcludedCities:    make(map[string]map[string]map[string]bool),
			},
		}
		err error
	)
	contractData := strings.Split(contractText, "\n")

	if len(contractData) < 2 {
		err = errors.New("Invalid contract, regions not found")
		return nil, err
	}
	heading := strings.TrimLeft(contractData[0], " ")
	if !strings.HasPrefix(heading, "Permissions for ") {
		err = errors.New("Invalid contract, heading line: Prefix: 'Permissions for ' not found")
		return nil, err
	}

	distributorHeirarchyText := strings.TrimPrefix(heading, "Permissions for ")
	distributorHeirarchyText = strings.ReplaceAll(distributorHeirarchyText, " ", "") //Remove spaces
	distributorHeirarchy := strings.Split(distributorHeirarchyText, "<")
	switch len(distributorHeirarchy) {
	case 0:
		return nil, errors.New("Invalid contract, distributor(s) not found in heading line after 'Permissions for': " + distributorHeirarchyText)
	case 1:
		if distributorHeirarchy[0] == "" {
			return nil, errors.New("Invalid contract, distributor(s) not found in heading line after 'Permissions for': " + distributorHeirarchyText)
		}
		contract.ContractRecipient = distributorHeirarchy[0]
	default:
		contract.ParentDistributor = &distributorHeirarchy[1]
		contract.ContractRecipient = distributorHeirarchy[0]
	}

	//check for duplication in distributor heirarchy, also check for empty strings
	distributorMap := make(map[string]bool)
	for _, distributor := range distributorHeirarchy {
		if distributor == "" {
			return nil, errors.New("Invalid contract, empty distributor found in heading line after 'Permissions for': " + distributorHeirarchyText)
		}
		if _, ok := distributorMap[distributor]; ok {
			return nil, errors.New("Invalid contract, duplicate distributor found in heading line after 'Permissions for': " + distributorHeirarchyText)
		}
		distributorMap[distributor] = true
	}

	for _, data := range contractData[1:] {
		data = strings.TrimLeft(data, " ")
		switch {
		case strings.HasPrefix(data, "INCLUDE: "):
			data = strings.TrimPrefix(data, "INCLUDE: ")
			err = contract.AddIncludedRegion(data)
			if err != nil {
				return nil, err
			}

		case strings.HasPrefix(data, "EXCLUDE: "):
			data = strings.TrimPrefix(data, "EXCLUDE: ")
			err = contract.AddExcludedRegion(data)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("Invalid contract, invalid line found: " + data)
		}
	}

	if len(contract.IncludedCountries) == 0 && len(contract.IncludedProvinces) == 0 && len(contract.IncludedCities) == 0 {
		err = errors.New("Invalid contract, no included regions found in contract")
		return nil, err
	}

	return &contract, nil
}

func (h *handler) GetDistributorPermissions(c *fiber.Ctx) error {
	distributor := c.Params("distributor")
	if distributor == "" {
		return response.InvalidURLParamResponse("distributor", errors.New("distributor not found in url")).WriteToJSON(c)
	}

	if c.Query("type", "text") == "json" {
		resp := h.databank.GetDistributorPermissionAsJSON(distributor)
		return resp.WriteToJSON(c)
	} else {
		note := h.databank.GetDistributorPermissionsAsText(distributor)
		return c.Status(200).SendString(note)
	}
}
