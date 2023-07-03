package cmd

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	_userController "github.com/ryanadiputraa/zenboard/app/user/controller"
	_userRepository "github.com/ryanadiputraa/zenboard/app/user/repository"
	_userService "github.com/ryanadiputraa/zenboard/app/user/service"

	_oauthController "github.com/ryanadiputraa/zenboard/app/oauth/controller"
	_oauthService "github.com/ryanadiputraa/zenboard/app/oauth/service"
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
	oauth := r.Group("/oauth")

	// user
	userRepository := _userRepository.NewUserRepository(DB)
	userService := _userService.NewUserService(userRepository)
	_userController.NewUserController(api, userService)

	// oauth
	oauthSerivce := _oauthService.NewOauthService()
	_oauthController.NewOauthController(oauth, oauthSerivce)

	r.Run(fmt.Sprintf(":%s", viper.GetString("PORT")))
}
