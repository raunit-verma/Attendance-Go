package adapter

import (
	"attendance/bean"

	"github.com/go-pg/pg"
)

func SetDBDev(cfg *bean.DBConfig) bean.DbDetails {
	return bean.DbDetails{
		User:     cfg.UserDev,
		Address:  cfg.AddressDev,
		Password: cfg.PasswordDev,
		Database: cfg.DatabaseDev,
	}
}

func SetDBProd(cfg *bean.DBConfig, addr string) bean.DbDetails {
	return bean.DbDetails{
		User:     cfg.UserProd,
		Address:  addr,
		Password: cfg.PasswordProd,
		Database: cfg.DatabaseProd,
	}
}

func GetPgOptions(dbConfig bean.DbDetails) *pg.Options {
	return &pg.Options{
		User:     dbConfig.User,
		Password: dbConfig.Password,
		Addr:     dbConfig.Address,
		Database: dbConfig.Database,
	}
}
