package handler

import (
	"challenge16/internal/regions"
	"challenge16/internal/response"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (h *handler) GetCountries(c *fiber.Ctx) error {
	resp := regions.GetCountries()
	return resp.WriteToJSON(c)
}

func (h *handler) GetProvincesInCountry(c *fiber.Ctx) error {
	countryCode := c.Params("countryCode")
	if countryCode == "" {
		return response.CreateError(400, URL_PARAM_MISSING, fmt.Errorf("Country code is required")).WriteToJSON(c)
	}
	resp := regions.GetProvincesInCountry(countryCode)
	return resp.WriteToJSON(c)
}

func (h *handler) GetCitiesInProvince(c *fiber.Ctx) error {
	countryCode := c.Params("countryCode")
	provinceCode := c.Params("provinceCode")
	if countryCode == "" || provinceCode == "" {
		return response.CreateError(400, URL_PARAM_MISSING, fmt.Errorf("Country code and province code are required")).WriteToJSON(c)
	}
	resp := regions.GetCitiesInProvince(countryCode, provinceCode)
	return resp.WriteToJSON(c)
}
