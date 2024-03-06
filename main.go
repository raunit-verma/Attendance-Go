package main

import (
	"attendance/bean"
	"attendance/repository"
	"net/http"

	"github.com/caarlos0/env"
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

	cfg := bean.MainConfig{}

	if err := env.Parse(&cfg); err != nil {
		zap.L().Error("Error loading env.", zap.Error(err))
	}

	db := repository.GetDB()
	defer db.Close()

	app := InitializeApp(db)

	r := app.NewMUXRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.Url},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST"},
	})

	handler := c.Handler(r)

	defer zap.L().Sync()

	zap.L().Info(`Server starting on Port ` + cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		zap.L().Fatal("HTTP server failed to start at Port "+cfg.Port, zap.Error((err)))
	}
}
