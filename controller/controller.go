package controller

import (
	"context"

	"github.com/PickHD/pickablog/config"
	"github.com/PickHD/pickablog/service"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type (
	// IHealthCheckController is an interface that has all the function to be implemented inside health check controller
	IHealthCheckController interface {
		HealthCheck(ctx *fiber.Ctx) error
	}
	
	// HealthCheckController is an app health check struct that consists of all the dependencies needed for health check controller
	HealthCheckController struct {
		Context context.Context
		Config *config.Configuration
		Logger *logrus.Logger
		HealthCheckSrv service.IHealthCheckService
	}

)