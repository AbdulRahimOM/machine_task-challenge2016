package main

import (
	"challenge16/internal/handler"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	app := fiber.New()
	app.Use(logger.New())

	handler := handler.NewHandler()

	//initialize the routes
	{
		distributor := app.Group("/distributor")
		{
			distributor.Post("/", handler.AddDistributor)
			distributor.Delete("/:distributor", handler.RemoveDistributor)
			distributor.Get("/", handler.GetDistributors)
		}

		permission := app.Group("/permission")
		{
			permission.Get("/check", handler.CheckIfDistributionIsAllowed)
			permission.Post("/allow", handler.AllowDistribution)
			permission.Post("/contract", handler.ApplyContract)
			permission.Post("/disallow", handler.DisallowDistribution)
			permission.Get("/:distributor", handler.GetDistributorPermissions)
		}

		regions := app.Group("/regions")
		{
			regions.Get("/countries", handler.GetCountries)
			regions.Get("/provinces/:countryCode", handler.GetProvincesInCountry)
			regions.Get("/cities/:countryCode/:provinceCode", handler.GetCitiesInProvince)
		}
	}

	err := app.Listen(fmt.Sprintf(":4010"))
	if err != nil {
		panic("Couldn't start the server. Error: " + err.Error())
	}
}
