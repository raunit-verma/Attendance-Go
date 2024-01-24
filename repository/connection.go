package repository

import (
	"log"
	"os"

	"github.com/go-pg/pg"
)

type DbConfig struct {
	User     string
	Password string
	Address  string
	Database string
}

var pgDb *pg.DB = nil

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

	var db *pg.DB = pg.Connect(opts)

	if db == nil {
		log.Printf("Database connection failed.\n")
		os.Exit(100)
	}

	log.Printf("Postgres Connected\n")
	return db
}
