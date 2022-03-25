package application

import (
	"github.com/PickHD/pickablog/controller"
	"github.com/PickHD/pickablog/repository"
	"github.com/PickHD/pickablog/service"
)

// Dependency can contain anything that will provide data for controller layer
type Dependency struct {
	HealthCheckController controller.IHealthCheckController
}

// SetupDependencyInjection is a function to set up dependencies 
func SetupDependencyInjection(app *App) *Dependency {
	return &Dependency{
		HealthCheckController: setupHealthCheckDependency(app),
	}
}

// setupHealthCheckDependency is a function to set up dependencies to be used inside health check controller layer
func setupHealthCheckDependency(app *App) *controller.HealthCheckController {
	healthCheckRepo := &repository.HealthCheckRepository{
		Context: app.Context,
		Config: app.Config,
		Logger: app.Logger,
		DB: app.DB,
	}

	healthCheckSrv := &service.HealthCheckService{
		Context: app.Context,
		Config: app.Config,
		Logger: app.Logger,
		HealthCheckRepo: healthCheckRepo,
	}

	healthCheckCtrl := &controller.HealthCheckController{
		Context: app.Context,
		Config: app.Config,
		Logger: app.Logger,
		HealthCheckSrv: healthCheckSrv,
	}

	return healthCheckCtrl
}