package handler

import (
	"challenge16/internal/regions"
	"challenge16/internal/response"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

const (
	INVALID_REGION = "INVALID_REGION"
)

func (h *handler) GetCountries(c *fiber.Ctx) error {
	countries := regions.GetCountries()
	return response.CreateSuccess(200, "SUCCESS", map[string]interface{}{
		"countries": countries,
	}).WriteToJSON(c)
}

func (h *handler) GetProvincesInCountry(c *fiber.Ctx) error {
	countryCode := c.Params("countryCode")
	if countryCode == "" {
		return response.CreateError(400, URL_PARAM_MISSING, fmt.Errorf("Country code is required")).WriteToJSON(c)
	}
	if !regions.CheckCountry(countryCode) {
		return response.CreateError(400, INVALID_REGION, fmt.Errorf("Invalid country code")).WriteToJSON(c)
	}

	provinces := regions.GetProvincesInCountry(countryCode)
	return response.CreateSuccess(200, "SUCCESS", map[string]interface{}{
		"provinces": provinces,
	}).WriteToJSON(c)
}

func (h *handler) GetCitiesInProvince(c *fiber.Ctx) error {
	countryCode := c.Params("countryCode")
	provinceCode := c.Params("provinceCode")
	if countryCode == "" || provinceCode == "" {
		return response.CreateError(400, URL_PARAM_MISSING, fmt.Errorf("Country code and province code are required")).WriteToJSON(c)
	}
	if !regions.CheckCountry(countryCode) {
		return response.CreateError(400, INVALID_REGION, fmt.Errorf("Invalid country code")).WriteToJSON(c)
	}
	if !regions.CheckProvince(countryCode, provinceCode) {
		return response.CreateError(400, INVALID_REGION, fmt.Errorf("Invalid province code")).WriteToJSON(c)
	}
	cities := regions.GetCitiesInProvince(countryCode, provinceCode)
	return response.CreateSuccess(200, "SUCCESS", map[string]interface{}{
		"cities": cities,
	}).WriteToJSON(c)
}
