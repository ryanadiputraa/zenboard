package main

import (
	_ "github.com/lib/pq"

	log "github.com/sirupsen/logrus"

	"github.com/ryanadiputraa/zenboard/config"
	"github.com/ryanadiputraa/zenboard/internal/server"
	"github.com/ryanadiputraa/zenboard/pkg/db"
	"github.com/ryanadiputraa/zenboard/pkg/logger"
)

func main() {
	logger.SetupLoger()
	conf, err := config.LoadConfig("./config/config")
	if err != nil {
		log.Fatal("fail to load config: ", err)
	}

	db, err := db.NewDBConn(conf.Database)
	if err != nil {
		log.Fatalf("fail to open db connection: %s", err)
	}

	s := server.NewServer(conf, db)
	if err = s.Run(); err != nil {
		log.Fatal("fail to run server: ", err)
	}
}
