package repository

import (
	"context"

	"github.com/PickHD/pickablog/config"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

type (
	// IHealthCheckRepository is an interface that has all the function to be implemented inside health check repository
	IHealthCheckRepository interface {
		HealthCheck() (bool,error)
	}

	// HealthCheckRepository is an app health check struct that consists of all the dependencies needed for health check repository
	HealthCheckRepository struct {
		Context context.Context
		Config *config.Configuration
		Logger *logrus.Logger
		DB *pgx.Conn
	}
)

// HealthCheck pinging the databases is OK or not
func (hcr *HealthCheckRepository) HealthCheck() (bool,error) {
	if err := hcr.DB.Ping(hcr.Context); err != nil {
		return false,err
	}

	return true,nil
}