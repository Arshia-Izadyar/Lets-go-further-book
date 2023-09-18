package db

import (
	"clean_api/src/config"
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

var PsqlDb *sql.DB

func InitDB(cfg *config.Config) (*sql.DB, error) {
	var err error
	PsqlDb, err = sql.Open("postgres", cfg.DB.Dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = PsqlDb.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	PsqlDb.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	PsqlDb.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	d, err := time.ParseDuration(cfg.DB.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	PsqlDb.SetConnMaxIdleTime(d)

	return PsqlDb, nil
}

func GetDB() *sql.DB {
	return PsqlDb
}

func CloseDB() error {
	err := PsqlDb.Close()
	return err
}
