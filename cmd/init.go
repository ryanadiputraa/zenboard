package cmd

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	db "github.com/ryanadiputraa/zenboard/pkg/db/sqlc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	DB *db.Queries
	Tx db.Tx
)

func init() {
	setupLoger()
	loadConfig()

	conn, err := sql.Open(viper.GetString("DB_DRIVER"), viper.GetString("DB_SOURCE"))
	if err != nil {
		log.Fatalf("fail to open db connection: %s", err)
	}
	DB = db.New(conn)
	Tx = db.NewTx(conn)
}

func setupLoger() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

func loadConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("fail to load config file: %s", err)
	}
}
