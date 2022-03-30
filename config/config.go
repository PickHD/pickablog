package config

import (
	"github.com/PickHD/pickablog/helper"
)

type (
	// Configuration is application configuration
	Configuration struct {
		Const *Constants
		Database *Database
		Redis *Redis
		Secret *Secret
	}

	// Constants is used To store configurable value, not the constant-constant value
	Constants struct {
		HTTPPort int
		ENV string
		GRedirectURL string
		OauthGoogleAPIURL string
	}

	// Database configuration
	Database struct {
		DBUser string
		DBPassword string
		DBHost string
		DBPort int
		DBName string
	}

	Redis struct {
		RDBHost string
		RDBPort int
		RDBExpire int
	}

	// Secret configuration
	Secret struct {
		GClientSecret string
		GClientID string
		JWTSecret string
	}
)

// LoadConfiguration...
func LoadConfiguration() *Configuration {
	return &Configuration{
		Const: loadConstants(),
		Database: loadDatabase(),
		Redis: loadRedis(),
		Secret: loadSecret(),
	}
}

// loadConstants...
func loadConstants() *Constants {
	return &Constants{
		HTTPPort: helper.GetEnvInt("PORT"),
		ENV: helper.GetEnvString("ENV"),
		GRedirectURL: helper.GetEnvString("G_REDIRECT_URL"),
		OauthGoogleAPIURL: helper.GetEnvString("OAUTH_G_API_URL"),
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

// loadRedis...
func loadRedis() *Redis {
	return &Redis{
		RDBHost: helper.GetEnvString("RDB_HOST"),
		RDBPort: helper.GetEnvInt("RDB_PORT"),
		RDBExpire: helper.GetEnvInt("RDB_EXPIRE"),
	}
}

// loadSecret...
func loadSecret() *Secret {
	return &Secret{
		GClientID: helper.GetEnvString("G_CLIENT_ID"),
		GClientSecret: helper.GetEnvString("G_CLIENT_SECRET"),
		JWTSecret: helper.GetEnvString("JWT_SECRET"),
	}
}