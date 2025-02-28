package handler

import (
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
	contract := string(c.Body())
	distributorHeirarchy, includeRegions, excludeRegions, err := getContractData(contract)
	if err != nil {
		return response.Response{
			HttpStatusCode: 400,
			ResponseCode:   "INVALID_CONTRACT",
			Error:          err,
		}.WriteToJSON(c)
	}
	resp := h.databank.ApplyContract(distributorHeirarchy, includeRegions, excludeRegions)
	return resp.WriteToJSON(c)
}

func getContractData(contract string) (distributorHeirarchy, includeRegions, excludeRegions []string, err error) {
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

	contractData := strings.Split(contract, "\n")

	if len(contractData) < 2 {
		err = errors.New("Invalid contract, regions not found")
		return
	}
	heading := strings.TrimLeft(contractData[0], " ")
	if !strings.HasPrefix(heading, "Permissions for ") {
		err = errors.New("Invalid contract, heading line: Prefix: 'Permissions for ' not found")
		return
	} else {
		data := strings.TrimPrefix(heading, "Permissions for ")
		data = strings.TrimRight(data, " ") //to avoid spaces at the end made by mistake
		distributorHeirarchy = strings.Split(data, " < ")
		if len(distributorHeirarchy) == 0 {
			err = errors.New("Invalid contract, distributor(s) not found in heading line after 'Permissions for': " + data)
			return
		} else if len(distributorHeirarchy) == 1 {
			if distributorHeirarchy[0] == "" {
				err = errors.New("Invalid contract, distributor(s) not found in heading line after 'Permissions for': " + data)
				return
			}
		}
	}

	for _, data := range contractData[1:] {
		data = strings.TrimLeft(data, " ")
		switch {
		case strings.HasPrefix(data, "INCLUDE: "):
			data = strings.TrimPrefix(data, "INCLUDE: ")
			includeRegions = append(includeRegions, data)
		case strings.HasPrefix(data, "EXCLUDE: "):
			data = strings.TrimPrefix(data, "EXCLUDE: ")
			excludeRegions = append(excludeRegions, data)
		default:
			err = errors.New("Invalid contract, invalid line found: " + data)
		}
	}

	if len(distributorHeirarchy) == 0 {
		err = errors.New("Invalid contract, distributor(s) not found")
		return
	}

	if len(includeRegions) == 0 && len(excludeRegions) == 0 {
		err = errors.New("Invalid contract, no permissions found")
		return
	}

	return
}
