package server

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/ryanadiputraa/zenboard/config"
	"github.com/ryanadiputraa/zenboard/internal/middleware"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	gin  *gin.Engine
	conf *config.Config
	db   *sqlx.DB
}

func NewServer(conf *config.Config, db *sqlx.DB) *Server {
	return &Server{
		gin:  gin.Default(),
		conf: conf,
		db:   db,
	}
}

func (s *Server) Run() error {
	if s.conf.Server.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	s.gin.SetTrustedProxies(nil)
	s.gin.Use(customRecovery())
	s.gin.Use(middleware.CORSMiddleware())

	s.MapHandlers()

	return s.gin.Run(s.conf.Server.Port)
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
