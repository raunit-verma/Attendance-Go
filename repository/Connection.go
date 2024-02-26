package repository

import (
	"attendance/bean"

	"github.com/caarlos0/env"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

type DbConfig struct {
	User     string
	Password string
	Address  string
	Database string
}

var pgDb *pg.DB = nil

func GetDB() *pg.DB {

	cfg := bean.DBConfig{}
	if err := env.Parse(&cfg); err != nil {
		zap.L().Error("Error Loading Env.", zap.Error(err))
	}

	if pgDb == nil {
		if cfg.Type == "Development" {
			DbConfig := DbConfig{
				User:     cfg.UserDev,
				Address:  cfg.AddressDev,
				Password: cfg.PasswordDev,
				Database: cfg.DatabaseDev,
			}
			pgDb = ConnectToDB(DbConfig)
			zap.L().Info("Connection to Development Database.")
		} else {
			pgUrl, _ := pg.ParseURL(cfg.AddressProd)
			DbConfig := DbConfig{
				User:     cfg.UserProd,
				Address:  pgUrl.Addr,
				Password: cfg.PasswordProd,
				Database: cfg.DatabaseProd,
			}
			pgDb = ConnectToDB(DbConfig)
			zap.L().Info("Connection to Production Database.")
		}
		_ = CreateSchema(pgDb, cfg)
	}
	return pgDb
}

func ConnectToDB(dbConfig DbConfig) *pg.DB {

	opts := &pg.Options{
		User:     dbConfig.User,
		Password: dbConfig.Password,
		Addr:     dbConfig.Address,
		Database: dbConfig.Database,
	}

	db := pg.Connect(opts)

	if db == nil {
		zap.L().Fatal("Database connection failed")
	} else {
		zap.L().Info("Database connected", zap.String("type", db.String()))
	}
	return db
}
