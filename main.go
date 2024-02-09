package main

import (
	"attendance/api/router"
	"attendance/repository"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
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
		Port: os.Getenv("PORT"),
	}

	r := router.NewMUXRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("URL")},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	db := repository.GetDB()
	defer db.Close()
	defer zap.L().Sync()

	zap.L().Info(`Server starting on Port ` + serverConfig.Port)

	if err := http.ListenAndServe(":"+serverConfig.Port, handler); err != nil {
		zap.L().Fatal("HTTP server failed to start at Port "+serverConfig.Port, zap.Error((err)))
	}
}
