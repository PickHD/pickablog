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
}

// SetupDependencyInjection is a function to set up dependencies 
func SetupDependencyInjection(app *App) *Dependency {
	return &Dependency{
		HealthCheckController: setupHealthCheckDependency(app),
		AuthController: setupAuthDependency(app),
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