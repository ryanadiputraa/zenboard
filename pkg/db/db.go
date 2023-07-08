package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/ryanadiputraa/zenboard/config"
)

const (
	maxOpenConns    = 60
	connMaxLifeTime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

func NewDBConn(conf config.Database) (db *sqlx.DB, err error) {
	db, err = sqlx.Connect(conf.Driver, conf.DSN)
	if err != nil {
		return
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifeTime)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime)

	return
}
