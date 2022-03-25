package infrastructure

import (
	"github.com/PickHD/pickablog/application"
	"github.com/PickHD/pickablog/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// ServeHTTP is function to start the apps infra in HTTP mode
func ServeHTTP(app *application.App) *fiber.App {
	//call setup router
	setupRouter(app)
	
	return app.Application
}

// setupRouter is function to manage all routings
func setupRouter(app *application.App) {
	var dep = application.SetupDependencyInjection(app)

	app.Application.Use(cors.New())
 	api := app.Application.Group("/api")

	v1 := api.Group("/v1", func (ctx *fiber.Ctx) error {
		ctx.Set("Version","v1")
		return ctx.Next()
	})

	v1.Get("/health",dep.HealthCheckController.HealthCheck)
	v1.Options("/health",helper.OptionsHandler)
}