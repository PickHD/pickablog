package infrastructure

import (
	"github.com/PickHD/pickablog/application"
	"github.com/PickHD/pickablog/helper"
	m "github.com/PickHD/pickablog/middleware"
	"github.com/gofiber/fiber/v2"
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

		v1.Options("/*",helper.OptionsHandler)
		v1.Get("/health",dep.HealthCheckController.HealthCheck)
		
		// AUTH SECTION
		{
			v1.Post("/auth/register",dep.AuthController.RegisterAuthor)
			v1.Get("/auth/google/login",dep.AuthController.GoogleLogin)
			v1.Get("/auth/google/callback",dep.AuthController.GoogleLoginCallback)
			v1.Post("/auth/login",dep.AuthController.Login)
		}

		// TAG SECTION
		{	
			v1.Post("/tag",m.ValidateJWTMiddleware,m.SuperAdminOnlyMiddleware,dep.TagController.CreateTag)
			v1.Get("/tag",dep.TagController.GetAllTag)
			v1.Put("/tag/:id",m.ValidateJWTMiddleware,m.SuperAdminOnlyMiddleware,dep.TagController.UpdateTag)
			v1.Delete("/tag/:id",m.ValidateJWTMiddleware,m.SuperAdminOnlyMiddleware,dep.TagController.DeleteTag)
		}

		// USER SECTION
		{
			v1.Get("/users",m.ValidateJWTMiddleware,m.SuperAdminOnlyMiddleware,dep.UserController.GetAllUser)
			v1.Get("/users/:id",m.ValidateJWTMiddleware,m.SuperAdminOnlyMiddleware,dep.UserController.GetUser)
			v1.Put("/users/:id",m.ValidateJWTMiddleware,m.SuperAdminOnlyMiddleware,dep.UserController.UpdateUser)
			v1.Delete("/users/:id",m.ValidateJWTMiddleware,m.SuperAdminOnlyMiddleware,dep.UserController.DeleteUser)
		}
	}
	
	// handler for route not found
	app.Application.Use(func(c *fiber.Ctx) error {
   	 	return helper.ResponseFormatter[any](c,fiber.StatusNotFound,nil,"Route not found",nil,nil)
	})
}