package application

import (
	"context"
	"fmt"
	"time"

	"github.com/PickHD/pickablog/config"
	"github.com/PickHD/pickablog/helper"


	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	pgx "github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// App is an wrapper application instance that contains application context, configuration, logger, etc
type App struct {
	Application *fiber.App
	Context context.Context
	Config *config.Configuration
	Logger *logrus.Logger
	DB *pgx.Conn
	GConfig *oauth2.Config
}

// SetupApplication is a function to create application instance
func SetupApplication(ctx context.Context) (*App, error) {
	var err error

	app := &App{}
	app.Context = context.TODO()
	app.Config = config.LoadConfiguration()
	if err != nil {
		return app, err
	}
	// custom log app with logrus
	logWithLogrus := logrus.New()
	logWithLogrus.Formatter = &logrus.JSONFormatter{}
	app.Logger = logWithLogrus

	// setup fiber in separated func
	app.Application = setupFiber(fiber.New())
	
	app.DB,err = pgx.Connect(context.Background(),fmt.Sprintf("postgres://%s:%s@%s:%d/%s",app.Config.Database.DBUser,app.Config.Database.DBPassword,app.Config.Database.DBHost,app.Config.Database.DBPort,app.Config.Database.DBName))
	if err != nil {
		app.Logger.Error("Failed connecting to databases, reason :%v",err)
		return app,err
	}

	app.Logger.Info("Success connecting to database...")

	app.GConfig = &oauth2.Config{
		RedirectURL: app.Config.Const.GRedirectURL,
		ClientID: app.Config.Secret.GKey,
		ClientSecret: app.Config.Secret.GSecret,
		Scopes: []string{"email","profile"},
		Endpoint: google.Endpoint,
	}

	return app,nil
}

// setupFiber is function separated for fiber configuration
func setupFiber(app *fiber.App) *fiber.App {
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowHeaders: "Authorization,access_token",
	}))
	app.Use(requestid.New())
	app.Use(limiter.New(limiter.Config{
		Expiration: 3 * time.Minute,
		LimitReached: func (ctx *fiber.Ctx) error {
			return helper.ResponseFormatter[any](ctx,fiber.StatusTooManyRequests,nil,"Too many requests, wait till 3 min",nil)
		},
	}))
	app.Use(func(c *fiber.Ctx) error {
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

	return app
}

// Close is a function to gracefully close the application
func (app *App) Close() {
	if app.DB != nil {
		app.DB.Close(context.Background())
	}

	app.Logger.Info("APP SUCCESSFULLY CLOSED")
}
