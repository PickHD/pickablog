package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"

	"github.com/PickHD/pickablog/application"
	"github.com/PickHD/pickablog/infrastructure"
	"github.com/joho/godotenv"
)

const (
	localServerMode = "local"
	httpServerMode = "http"
)

func main() {
	_ = godotenv.Load("./config/dev.env")
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Checking command arguments
	var (
		args = os.Args[1:]
		mode = localServerMode
	)

	if len(args) > 0 {
		mode = os.Args[1]
	}

	// create a context with background for setup the application
	ctx := context.Background()
	app, err := application.SetupApplication(ctx)
	if err != nil {
		app.Logger.Error("Failed to initialize app. Error: ", err)
		panic(err)
	}

	//create a channel for listening to OS signals and connecting OS interrupts to the channel
	c := make(chan os.Signal, 1)
	signal.Notify(c,os.Interrupt)
	serverShutdown := make(chan struct{})

	go func() {
		_ = <-c
		app.Logger.Info("APP GRACEFULLY SHUTDOWN")
		app.Close()
		_ = app.Application.Shutdown()
		serverShutdown <- struct{}{}
	}()

	switch mode {
	case localServerMode,httpServerMode:
		var (
			httpServer = infrastructure.ServeHTTP(app)
		)

		if err := httpServer.Listen(fmt.Sprintf(":%d",app.Config.Const.HTTPPort)); err != nil {
		 	app.Logger.Error("Failed to run app. Error: ", err)
		    	panic(err)
		}

	}
}
