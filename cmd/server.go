package cmd

import (
	"fmt"

	"github.com/gin-gonic/gin"
	_userController "github.com/ryanadiputraa/zenboard/app/user/controller"
	_userRepository "github.com/ryanadiputraa/zenboard/app/user/repository"
	_userService "github.com/ryanadiputraa/zenboard/app/user/service"
	"github.com/spf13/viper"
)

func ServeHTTP() {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	if viper.GetString("ENV") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(CORSMiddleware())

	api := r.Group("/api")

	// user
	userRepository := _userRepository.NewUserRepository(DB)
	userService := _userService.NewUserService(userRepository)
	_userController.NewUserController(api, userService)

	r.Run(fmt.Sprintf(":%s", viper.GetString("PORT")))
}
