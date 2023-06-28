package cmd

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func ServeHTTP() {
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(":%s", viper.GetString("PORT")), nil),
	)
}
