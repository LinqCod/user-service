package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/linqcod/user-service/cmd/api/route"
	"github.com/linqcod/user-service/pkg/config"
	"github.com/linqcod/user-service/pkg/database"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	config.LoadConfig(".env")
}

func main() {
	// init zap logger
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	baseLogger, err := loggerConfig.Build()
	if err != nil {
		log.Fatalf("error while building zap logger: %v", err)
	}

	logger := baseLogger.Sugar()

	// init db connection
	username := viper.GetString("POSTGRES_USER")
	password := viper.GetString("POSTGRES_PASSWORD")
	port := viper.GetString("POSTGRES_PORT")
	dbname := viper.GetString("POSTGRES_DB")

	db, err := database.InitDB(username, password, port, dbname)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Fatalf("error while trying to ping db: %v", err)
	}

	// init routing
	contextTimeout := viper.GetInt64("CONTEXT_TIMEOUT")
	timeout := time.Duration(contextTimeout) * time.Second
	validate := validator.New()
	router := gin.Default()

	route.Setup(logger, timeout, db, validate, router)

	// init server
	serverAddr := fmt.Sprintf(":%s", viper.GetString("SERVER_PORT"))
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// graceful shutdown
	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Fatalf("error while trying to shutdown http server: %v", err)
		}
		close(stopped)
	}()

	logger.Infof("Starting HTTP server on %s", serverAddr)

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-stopped

	log.Printf("Have a nice day :)")
}
