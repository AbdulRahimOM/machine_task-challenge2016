package handler

import (
	"challenge16/validation"

	"github.com/gofiber/fiber/v2"
)

type SS struct {
	Distributor  string `query:"distributor" validate:"required"`
	RegionString string `query:"region" validate:"required"`
}

func (h *handler) CheckIfDistributionIsAllowed(c *fiber.Ctx) error {
	req := new(SS)

	if ok, err := validation.BindAndValidateURLQueryRequest(c, req); !ok {
		return err
	}

	resp := h.databank.IsAllowed(req.Distributor, req.RegionString)
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
