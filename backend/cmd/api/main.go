package main

import (
	"context"
	"fmt"
	"hris-backend/internal/bootstrap"
	"hris-backend/internal/modules/attendance"
	"hris-backend/internal/routes"
	"hris-backend/internal/seeder"
	"hris-backend/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/joho/godotenv"
)

const (
	httpServerMode = "http"
	seedMode       = "seed"
)

// @title HRIS Backend
// @version 1.0.0
// @description A REST API built with Go, Echo, and Clean Architecture
// @host localhost:8080
// @BasePath /api/v1
func main() {
	time.Local = time.UTC
	runtime.GOMAXPROCS(runtime.NumCPU())
	args := os.Args[1:]
	mode := httpServerMode
	if len(args) > 0 {
		mode = args[0]
	}

	envPaths := []string{
		"../../../.env", "../../.env", "../.env", ".env", "/app/.env",
	}

	var loadErr error
	for _, path := range envPaths {
		if loadErr = godotenv.Load(path); loadErr == nil {
			break
		}
	}

	if loadErr != nil {
		logger.Warn("Warning: .env file not found (this is OK in Docker, env vars will be used)")
	}

	appContainer, err := bootstrap.NewContainer()
	if err != nil {
		logger.Errorw("Failed to initialize application container: ", err)
		os.Exit(1)
	}
	defer appContainer.Close()

	switch mode {
	case seedMode:
		if err := seeder.Execute(appContainer.DB.GetDB(), appContainer.Config, appContainer.Bcrypt); err != nil {
			logger.Errorw("Seeding failed: ", err)
			os.Exit(1)
		}
	case httpServerMode:
		attendance.StartGeocodeWorker(appContainer.DB.GetDB(), appContainer.Location)

		logger.Info("Starting HRIS API Server...")

		appRouter := routes.ServeHTTP(appContainer)

		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", appContainer.Config.Server.Port),
			Handler: appRouter,
		}

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)

		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Errorw("Failed to to start server. Error: ", err)
			}
		}()

		<-sigCh

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Errorw("Failed to shutdown server. Error: ", err)
		}

		logger.Info("Server Shutdown Gracefully...")
	}

}
