package db

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/ryanadiputraa/zenboard/config"
)

const (
	maxOpenConns    = 60
	connMaxLifeTime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

func NewConn(conf *config.Config) (db *sql.DB, err error) {
	db, err = sql.Open(conf.Database.Driver, conf.Database.DSN)
	if err != nil {
		return
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifeTime)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime)

	err = db.Ping()
	return
}
