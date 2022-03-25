package application

import (
	"context"
	"fmt"

	"github.com/PickHD/pickablog/config"

	"github.com/gofiber/fiber/v2"
	pgx "github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

// App is an wrapper application instance that contains application context, configuration, logger, etc
type App struct {
	Application *fiber.App
	Context context.Context
	Config *config.Configuration
	Logger *logrus.Logger
	DB *pgx.Conn
}

// SetupApplication is a function to create application instance
func SetupApplication(ctx context.Context) (*App, error) {
	var err error

	app := &App{}
	app.Application = fiber.New()
	app.Context = context.TODO()
	app.Config = config.LoadConfiguration()
	app.Logger = logrus.New()
	if err != nil {
		return app, err
	}

	app.Application.Use(func(c *fiber.Ctx) error {
		// Set some security headers:
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Download-Options", "noopen")
		c.Set("Strict-Transport-Security", "max-age=5184000")
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("X-DNS-Prefetch-Control", "off")

		// Go to next middleware:
		return c.Next()
	})

	app.DB,err = pgx.Connect(context.Background(),fmt.Sprintf("postgres://%s:%s@%s:%d/%s",app.Config.Database.DBUser,app.Config.Database.DBPassword,app.Config.Database.DBHost,app.Config.Database.DBPort,app.Config.Database.DBName))
	if err != nil {
		app.Logger.Error("Failed connecting to databases, reason :%v",err)
		return app,err
	}

	app.Logger.Info("Success connecting to database...")

	//TODO : instantiate google oauth here

	return app,nil
}

// Close is a function to gracefully close the application
func (app *App) Close() {
	if app.DB != nil {
		app.DB.Close(context.Background())
	}

	app.Logger.Info("APP SUCCESSFULLY CLOSED")
}
