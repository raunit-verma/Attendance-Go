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
		zap.L().Error("Error loading Env", zap.Error(err))
	}

	serverConfig := router.ServerConfig{
		Port: os.Getenv("PORT"),
	}
	db := repository.GetDB()
	defer db.Close()

	wire := InitializeApp(db)

	r := wire.NewMUXRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("URL")},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST"},
	})

	handler := c.Handler(r)

	defer zap.L().Sync()

	zap.L().Info(`Server starting on Port ` + serverConfig.Port)

	if err := http.ListenAndServe(":"+serverConfig.Port, handler); err != nil {
		zap.L().Fatal("HTTP server failed to start at Port "+serverConfig.Port, zap.Error((err)))
	}
}
