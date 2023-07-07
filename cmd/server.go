package cmd

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	_userController "github.com/ryanadiputraa/zenboard/app/user/controller"
	_userRepository "github.com/ryanadiputraa/zenboard/app/user/repository"
	_userService "github.com/ryanadiputraa/zenboard/app/user/service"

	_boardRepository "github.com/ryanadiputraa/zenboard/app/board/repository"
	_boardService "github.com/ryanadiputraa/zenboard/app/board/service"

	_oauthController "github.com/ryanadiputraa/zenboard/app/oauth/controller"
	_oauthService "github.com/ryanadiputraa/zenboard/app/oauth/service"
)

func ServeHTTP() {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	if viper.GetString("ENV") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(customRecovery())
	r.Use(gin.Logger())
	r.Use(CORSMiddleware())

	api := r.Group("/api")
	oauth := r.Group("/oauth")

	// user
	userRepository := _userRepository.NewUserRepository(DB)
	userService := _userService.NewUserService(userRepository)
	_userController.NewUserController(api, userService)

	// board
	boardRepository := _boardRepository.NewBoardRepository(DB, Tx)
	boardService := _boardService.NewBoardService(boardRepository)

	// oauth
	oauthSerivce := _oauthService.NewOauthService()
	_oauthController.NewOauthController(oauth, oauthSerivce, userService, boardService)

	r.Run(fmt.Sprintf(":%s", viper.GetString("PORT")))
}

func customRecovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error(err)
				debug.PrintStack()
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "something went wrong, please try again later",
				})
			}
		}()

		ctx.Next()
	}
}
