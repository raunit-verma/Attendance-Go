package repository

import (
	"os"

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

// remove use of contansts and use utils ( bin )
func GetDB() *pg.DB {
	DbConfig := DbConfig{
		User:     os.Getenv("DB_USER"),
		Address:  os.Getenv("DB_ADDRESS"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_DATABASE"),
	}
	if pgDb == nil {
		pgDb = ConnectToDB(DbConfig)
		_ = CreateSchema(pgDb)
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
