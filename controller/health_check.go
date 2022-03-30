package controller

import (
	"context"

	"github.com/PickHD/pickablog/config"
	"github.com/PickHD/pickablog/helper"
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
		HealthCheckSvc service.IHealthCheckService
	}

)

// HealthCheck controller layer to checking databases is ok or not
func (hcc *HealthCheckController) HealthCheck(ctx *fiber.Ctx) error {
	ok,err := hcc.HealthCheckSvc.HealthCheck()
	if err != nil || !ok {
		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,"Failed checking health services",nil)
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusOK,nil,"OK",nil)
} 