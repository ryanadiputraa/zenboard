package cmd

import (
	"fmt"
	"net/http"

	_userRepository "github.com/ryanadiputraa/zenboard/app/user/repository"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func ServeHTTP() {
	_userRepository.NewUserRepository(DB)

	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(":%s", viper.GetString("PORT")), nil),
	)
}
