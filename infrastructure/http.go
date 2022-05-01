package infrastructure

import (
	"time"

	"github.com/PickHD/pickablog/application"
	"github.com/PickHD/pickablog/helper"
	m "github.com/PickHD/pickablog/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
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
			v1.Post("/auth/login",limiter.New(limiter.Config{
				Expiration: 15 * time.Minute,
				Max: 5,
				LimitReached: func (ctx *fiber.Ctx) error {
					return helper.ResponseFormatter[any](ctx,fiber.StatusTooManyRequests,nil,"Login Attempts already reached the limit, tell our super admin about resetting a password, Thank you",nil,nil)
				},
			}),dep.AuthController.Login)
		}

		// TAG SECTION
		{	
			v1.Post("/tag",m.ValidateJWTMiddleware,m.SuperAdminOnlyMiddleware,dep.TagController.CreateTag)
			v1.Get("/tag",dep.TagController.ListTag)
			v1.Put("/tag/:id",m.ValidateJWTMiddleware,m.SuperAdminOnlyMiddleware,dep.TagController.UpdateTag)
			v1.Delete("/tag/:id",m.ValidateJWTMiddleware,m.SuperAdminOnlyMiddleware,dep.TagController.DeleteTag)
		}

		// USER SECTION
		{
			v1.Get("/users",m.ValidateJWTMiddleware,m.SuperAdminOnlyMiddleware,dep.UserController.ListUser)
			v1.Get("/users/:id",m.ValidateJWTMiddleware,m.SuperAdminOnlyMiddleware,dep.UserController.DetailUser)
			v1.Put("/users/:id",m.ValidateJWTMiddleware,m.SuperAdminOnlyMiddleware,dep.UserController.UpdateUser)
			v1.Delete("/users/:id",m.ValidateJWTMiddleware,m.SuperAdminOnlyMiddleware,dep.UserController.DeleteUser)
		}

		// BLOG SECTION
		{
			v1.Post("/blog",m.ValidateJWTMiddleware,m.AuthorOnlyMiddleware,dep.BlogController.CreateBlog)
			v1.Get("/blog",dep.BlogController.ListBlog)
			v1.Get("/blog/:slug",dep.BlogController.DetailBlog)
			v1.Put("/blog/:id",m.ValidateJWTMiddleware,m.AuthorOnlyMiddleware,dep.BlogController.UpdateBlog)
			v1.Delete("/blog/:id",m.ValidateJWTMiddleware,m.AuthorOnlyMiddleware,dep.BlogController.DeleteBlog)

			// BLOG COMMENT SECTION
			v1.Post("/blog/:id/comment",m.ValidateJWTMiddleware,dep.BlogController.CreateComment)
			v1.Put("/blog/:id/comment/:comment_id",m.ValidateJWTMiddleware,dep.BlogController.UpdateComment)
			v1.Get("/blog/:id/comment",m.ValidateJWTMiddleware,dep.BlogController.ListComment)
			v1.Delete("/blog/:id/comment/:comment_id",m.ValidateJWTMiddleware,dep.BlogController.DeleteComment)

			// BLOG LIKE SECTION
			v1.Get("/blog/:id/like",m.ValidateJWTMiddleware,dep.BlogController.Like)
			v1.Delete("/blog/:id/like/:like_id",m.ValidateJWTMiddleware,dep.BlogController.UnLike)
		}
	}
	
	// handler for route not found
	app.Application.Use(func(c *fiber.Ctx) error {
   	 	return helper.ResponseFormatter[any](c,fiber.StatusNotFound,nil,"Route not found",nil,nil)
	})
}