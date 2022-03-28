package infrastructure

import (
	"github.com/PickHD/pickablog/application"
	"github.com/PickHD/pickablog/helper"
	appware "github.com/PickHD/pickablog/middleware"

	"github.com/gofiber/fiber/v2"
	_ "github.com/gofiber/jwt/v3"

)

// ServeHTTP is wrapper function to start the apps infra in HTTP mode
func ServeHTTP(app *application.App) *fiber.App {
	//call setup router
	setupRouter(app)
	
	return app.Application
}

// setupRouter is function to manage all routings
func setupRouter(app *application.App) {
	var dep = application.SetupDependencyInjection(app)

 	api := app.Application.Group("/api")
	{
		v1 := api.Group("/v1", func (ctx *fiber.Ctx) error {
			ctx.Set("Version","v1")
			return ctx.Next()
		})
		
		{
			v1.Get("/health",dep.HealthCheckController.HealthCheck)
			v1.Options("/health",helper.OptionsHandler)

			// UNPROTECTED ROUTE HERE
			v1.Get("/test-protecc",appware.Chain(dep.HealthCheckController.HealthCheck,appware.ValidateJWTMiddleware))
			
			// v1.Use(jwtware.New(jwtware.Config{
			// 	SigningKey: []byte(app.Config.Secret.JWTSecret),
			// }))

			// PROTECTED ROUTE HERE

			
		}

	
	}
	
	// handler for route not found
	app.Application.Use(func(c *fiber.Ctx) error {
   	 	return helper.ResponseFormatter[any](c,404,nil,"Route not found",nil)
	})
}