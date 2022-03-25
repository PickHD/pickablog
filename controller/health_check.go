package controller

import (
	"github.com/PickHD/pickablog/helper"
	"github.com/gofiber/fiber/v2"
)

func (hcc *HealthCheckController) HealthCheck(ctx *fiber.Ctx) error {
	ok,err := hcc.HealthCheckSrv.HealthCheck()
	if err != nil || !ok {
		return helper.ResponseFormatter[any](ctx,fiber.StatusInternalServerError,err,"Failed checking health services",nil)
	}

	return helper.ResponseFormatter[any](ctx,fiber.StatusOK,nil,"OK",nil)
} 