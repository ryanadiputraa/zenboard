package server

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/ryanadiputraa/zenboard/config"
	"github.com/ryanadiputraa/zenboard/internal/middleware"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

type Server struct {
	gin  *gin.Engine
	conf *config.Config
	db   *sqlx.DB
	ws   *WebSocketServer
}

func NewServer(conf *config.Config, db *sqlx.DB) *Server {
	ws := &WebSocketServer{
		conns: make(map[string]map[*websocket.Conn]bool),
	}
	return &Server{
		gin:  gin.Default(),
		conf: conf,
		db:   db,
		ws:   ws,
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
