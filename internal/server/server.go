package server

import (
	"challenge16/internal/handler"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func NewServer(rateLimit int) *fiber.App {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(limiter.New(limiter.Config{
		Max:        rateLimit,
		Expiration: 1 * time.Minute,
	}))

	handler := handler.NewHandler()

	// Initialize the routes
	{
		// Distributor routes
		distributor := app.Group("/distributor")
		{
			distributor.Post("/", handler.AddDistributor)
			distributor.Delete("/:distributor", handler.RemoveDistributor)
			distributor.Get("/", handler.GetDistributors)
		}

		// Permission routes
		permission := app.Group("/permission")
		{
			permission.Get("/check", handler.CheckIfDistributionIsAllowed)
			permission.Post("/allow", handler.AllowDistribution)
			permission.Post("/contract", handler.ApplyContract)
			permission.Post("/disallow", handler.DisallowDistribution)
			permission.Get("/:distributor", handler.GetDistributorPermissions)
		}

		// Region routes
		regions := app.Group("/regions")
		{
			regions.Get("/countries", handler.GetCountries)
			regions.Get("/provinces/:countryCode", handler.GetProvincesInCountry)
			regions.Get("/cities/:countryCode/:provinceCode", handler.GetCitiesInProvince)
		}
	}

	return app
}
