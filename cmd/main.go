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
			distributor.Post("/add-sub", handler.AddSubDistributor)
		}

		permission := app.Group("/permission")
		{
			permission.Get("/check", handler.CheckIfDistributionIsAllowed)
			permission.Post("/allow", handler.AllowDistribution)
			permission.Post("/disallow", handler.DisallowDistribution)
		}
	}

	err := app.Listen(fmt.Sprintf(":4010"))
	if err != nil {
		panic("Couldn't start the server. Error: " + err.Error())
	}
}
