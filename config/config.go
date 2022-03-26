package config

import (
	"github.com/PickHD/pickablog/helper"
)

type (
	// Configuration is application configuration
	Configuration struct {
		Const *Constants
		Database *Database
		Secret *Secret
	}

	// Constants is used To store configurable value, not the constant-constant value
	Constants struct {
		HTTPPort int
		ENV string
	}

	// Database configuration
	Database struct {
		DBUser string
		DBPassword string
		DBHost string
		DBPort int
		DBName string
	}

	// Secret configuration
	Secret struct {
		GSecret string
		GKey string
	}
)

// LoadConfiguration...
func LoadConfiguration() *Configuration {
	return &Configuration{
		Const: loadConstants(),
		Database: loadDatabase(),
		Secret: loadSecret(),
	}
}

// loadConstants...
func loadConstants() *Constants {
	return &Constants{
		HTTPPort: helper.GetEnvInt("PORT"),
		ENV: helper.GetEnvString("ENV"),
	}
}

// loadDatabase...
func loadDatabase() *Database {
	return &Database{
		DBHost: helper.GetEnvString("DB_HOST"),
		DBUser:helper.GetEnvString("DB_USER"),
		DBPassword: helper.GetEnvString("DB_PASSWORD"),
		DBPort: helper.GetEnvInt("DB_PORT"),
		DBName: helper.GetEnvString("DB_NAME"),
	}
}

// loadSecret...
func loadSecret() *Secret {
	return &Secret{
		GKey: helper.GetEnvString("G_KEY"),
		GSecret: helper.GetEnvString("G_SECRET"),
	}
}