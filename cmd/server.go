package cmd

import (
	"database/sql"
	"net/http"
	"runtime/debug"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/ryanadiputraa/zenboard/config"
	db "github.com/ryanadiputraa/zenboard/pkg/db/sqlc"
	"github.com/ryanadiputraa/zenboard/pkg/logger"

	_userController "github.com/ryanadiputraa/zenboard/app/user/controller"
	_userRepository "github.com/ryanadiputraa/zenboard/app/user/repository"
	_userService "github.com/ryanadiputraa/zenboard/app/user/service"

	_boardRepository "github.com/ryanadiputraa/zenboard/app/board/repository"
	_boardService "github.com/ryanadiputraa/zenboard/app/board/service"

	_oauthController "github.com/ryanadiputraa/zenboard/app/oauth/controller"
	_oauthService "github.com/ryanadiputraa/zenboard/app/oauth/service"
)

func ServeHTTP() {
	logger.SetupLoger()
	conf, err := config.LoadConfig("./config/config")
	if err != nil {
		log.Fatal("fail to load config: ", err)
	}

	conn, err := sql.Open(conf.Database.Driver, conf.Database.DSN)
	if err != nil {
		log.Fatalf("fail to open db connection: %s", err)
	}

	DB := db.New(conn)
	Tx := db.NewTx(conn)

	r := gin.Default()
	r.SetTrustedProxies(nil)

	if conf.Server.Env == "prod" {
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
	_userController.NewUserController(conf, api, userService)

	// board
	boardRepository := _boardRepository.NewBoardRepository(DB, Tx)
	boardService := _boardService.NewBoardService(boardRepository)

	// oauth
	oauthSerivce := _oauthService.NewOauthService(conf)
	_oauthController.NewOauthController(conf, oauth, oauthSerivce, userService, boardService)

	r.Run(conf.Server.Port)
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
