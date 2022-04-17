package application

import (
	"github.com/PickHD/pickablog/controller"
	"github.com/PickHD/pickablog/repository"
	"github.com/PickHD/pickablog/requester"
	"github.com/PickHD/pickablog/service"
)

// Dependency can contain anything that will provide data for controller layer
type Dependency struct {
	HealthCheckController controller.IHealthCheckController
	AuthController controller.IAuthController
	TagController controller.ITagController
	UserController controller.IUserController
	BlogController controller.IBlogController
}

// SetupDependencyInjection is a function to set up dependencies 
func SetupDependencyInjection(app *App) *Dependency {
	return &Dependency{
		HealthCheckController: setupHealthCheckDependency(app),
		AuthController: setupAuthDependency(app),
		TagController: setupTagDependency(app),
		UserController: setupUserDependency(app),
		BlogController: setupBlogDependency(app),
	}
}

// setupHealthCheckDependency is a function to set up dependencies to be used inside health check controller layer
func setupHealthCheckDependency(app *App) *controller.HealthCheckController {
	healthCheckRepo := &repository.HealthCheckRepository{
		Context: app.Context,
		Config: app.Config,
		Logger: app.Logger,
		DB: app.DB,
		Redis: app.Redis,
	}

	healthCheckSvc := &service.HealthCheckService{
		Context: app.Context,
		Config: app.Config,
		Logger: app.Logger,
		HealthCheckRepo: healthCheckRepo,
	}

	healthCheckCtrl := &controller.HealthCheckController{
		Context: app.Context,
		Config: app.Config,
		Logger: app.Logger,
		HealthCheckSvc: healthCheckSvc,
	}

	return healthCheckCtrl
}

// setupAuthDependency is a function to set up dependencies to be used inside auth controller layer
func setupAuthDependency(app *App) *controller.AuthController {
	// init requester
	GOAuthReq := &requester.OAuthGoogle{
		Context: app.Context,
		Config: app.Config,
		Logger: app.Logger,
		GConfig: app.GConfig,
		HTTPClient: app.HTTPClient,
	}

	authRepo := &repository.AuthRepository{
		Context: app.Context,
		Config: app.Config,
		Logger: app.Logger,
		DB: app.DB,
		Redis: app.Redis,
	}

	authSvc := &service.AuthService{
		Context: app.Context,
		Config: app.Config,
		Logger: app.Logger,
		AuthRepo: authRepo,
		GConfig: app.GConfig,
		GOAuthReq: GOAuthReq,
	}

	authCtrl := &controller.AuthController{
		Context: app.Context,
		Config: app.Config,
		Logger: app.Logger,
		AuthSvc: authSvc,
	}

	return authCtrl
}

// setupTagDependency is a function to set up dependencies to be used inside tag controller layer
func setupTagDependency(app *App) *controller.TagController {
	tagRepo := &repository.TagRepository{
		Context: app.Context,
		Config: app.Config,
		Logger: app.Logger,
		DB: app.DB,
	}

	tagSvc := &service.TagService{
		Context: app.Context,
		Config: app.Config,
		Logger: app.Logger,
		TagRepo: tagRepo,
	}

	tagCtrl := &controller.TagController{
		Context: app.Context,
		Config: app.Config,
		Logger: app.Logger,
		TagSvc: tagSvc,
	}

	return tagCtrl
}

// setupUserDependency is a function to set up dependencies to be used inside user controller layer
func setupUserDependency(app *App) *controller.UserController {
	userRepo := &repository.UserRepository{
		Context: app.Context,
		Config: app.Config,
		Logger: app.Logger,
		DB: app.DB,
	}

	userSvc := &service.UserService{
		Context: app.Context,
		Config: app.Config,
		Logger: app.Logger,
		UserRepo: userRepo,
	}

	userCtrl := &controller.UserController{
		Context: app.Context,
		Config: app.Config,
		Logger: app.Logger,
		UserSvc: userSvc,
	}

	return userCtrl
}

// setupBlogDependency is a function to set up dependencies to be used inside blog controller layer
func setupBlogDependency(app *App) *controller.BlogController {
	blogRepo := &repository.BlogRepository{
		Context: app.Context,
		Config: app.Config,
		Logger: app.Logger,
		DB: app.DB,
	}

	commentRepo := &repository.CommentRepository{
		Context: app.Context,
		Config: app.Config,
		Logger: app.Logger,
		DB: app.DB,
	}

	likeRepo := &repository.LikeRepository{
		Context: app.Context,
		Config: app.Config,
		Logger: app.Logger,
		DB: app.DB,
	}

	blogSvc := &service.BlogService{
		Context: app.Context,
		Config: app.Config,
		Logger: app.Logger,
		BlogRepo: blogRepo,
		CommentRepo: commentRepo,
		LikeRepo: likeRepo,
	}

	blogCtrl := &controller.BlogController{
		Context: app.Context,
		Config: app.Config,
		Logger: app.Logger,
		BlogSvc: blogSvc,
	}

	return blogCtrl
}