package main

import (
	"attendance/api/router"
	"attendance/repository"
	"net/http"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func main() {
	err := godotenv.Load()

	if err != nil {
		zap.L().Fatal("Error loading Env")
	}

	serverConfig := router.ServerConfig{
		Port: ":1025",
	}

	r := router.NewMUXRouter()

	db := repository.GetDB()
	defer db.Close()
	defer zap.L().Sync()

	zap.L().Info(`Server starting on Port ` + serverConfig.Port)

	if err := http.ListenAndServe(serverConfig.Port, r); err != nil {
		zap.L().Fatal("HTTP server failed to start at Port "+serverConfig.Port, zap.Error((err)))
	}
}
