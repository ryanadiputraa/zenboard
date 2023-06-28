package cmd

import (
	"fmt"
	"net/http"

	_userRepository "github.com/ryanadiputraa/zenboard/app/user/repository"
	_userService "github.com/ryanadiputraa/zenboard/app/user/service"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func ServeHTTP() {
	userRepository := _userRepository.NewUserRepository(DB)
	_userService.NewUserService(userRepository)

	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(":%s", viper.GetString("PORT")), nil),
	)
}
