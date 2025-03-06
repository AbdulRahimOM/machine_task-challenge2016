package handler

import (
	"challenge16/internal/response"
	"challenge16/utils/validation"
	"errors"

	"github.com/gofiber/fiber/v2"
)

func (h *handler) AddDistributor(c *fiber.Ctx) error {
	req := new(struct {
		Distributor string `json:"distributor" validate:"required"`
	})

	if ok, err := validation.BindAndValidateJSONRequest(c, req); !ok {
		return err
	}

	resp := h.databank.AddDistributor(req.Distributor)
	return resp.WriteToJSON(c)
}

func (h *handler) RemoveDistributor(c *fiber.Ctx) error {
	distributor := c.Params("distributor")
	if distributor == "" {
		return response.CreateError(400, URL_PARAM_MISSING, errors.New("distributor is required")).WriteToJSON(c)
	}

	resp := h.databank.RemoveDistributor(distributor)
	return resp.WriteToJSON(c)
}

func (h *handler) GetDistributors(c *fiber.Ctx) error {
	resp := h.databank.GetDistributors()
	return resp.WriteToJSON(c)
}
