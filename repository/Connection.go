package repository

import (
	"attendance/adapter"
	"attendance/bean"
	"attendance/util"

	"github.com/caarlos0/env"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

var pgDb *pg.DB = nil

func GetDB() *pg.DB {

	cfg := bean.DBConfig{}
	if err := env.Parse(&cfg); err != nil {
		zap.L().Error("Error Loading Env.", zap.Error(err))
	}

	if pgDb == nil {
		if cfg.Type == util.DEVELOPMENT {
			pgDb = ConnectToDB(adapter.SetDBDev(&cfg))
			zap.L().Info("Connection to Development Database.")
		} else {
			pgUrl, _ := pg.ParseURL(cfg.AddressProd)
			pgDb = ConnectToDB(adapter.SetDBProd(&cfg, pgUrl.Addr))
			zap.L().Info("Connection to Production Database.")
		}
		_ = CreateSchema(pgDb, cfg)
	}
	return pgDb
}

func ConnectToDB(dbConfig bean.DbDetails) *pg.DB {

	db := pg.Connect(adapter.GetPgOptions(dbConfig))

	if db == nil {
		zap.L().Fatal("Database connection failed")
	} else {
		zap.L().Info("Database connected", zap.String("type", db.String()))
	}
	return db
}
