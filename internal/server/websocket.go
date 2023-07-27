package server

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/ryanadiputraa/zenboard/config"
	"github.com/ryanadiputraa/zenboard/internal/domain"
	"golang.org/x/net/websocket"

	log "github.com/sirupsen/logrus"
)

type WebSocketServer struct {
	conns map[string]map[*websocket.Conn]bool
	sync.Mutex
}

type socket struct {
	ctx    *gin.Context
	conf   config.JWT
	conn   *websocket.Conn
	roomID string
}

type wsService struct {
	boardService domain.BoardService
}

type webSocketMessage struct {
	Key  string `json:"key"`
	Data any    `json:"data"`
}

func (ws *WebSocketServer) HandleConnection(
	ctx *gin.Context,
	conf config.JWT,
	service wsService,
) {
	websocket.Handler(func(c *websocket.Conn) {
		boardID := ctx.Query("board_id")

		ws.Lock()
		_, ok := ws.conns[boardID]
		if !ok {
			ws.conns[boardID] = map[*websocket.Conn]bool{
				c: true,
			}
		} else {
			ws.conns[boardID][c] = true
		}
		ws.Unlock()
		log.Info(fmt.Sprintf("new connection on (%v) : %v", boardID, c))

		socket := &socket{
			ctx:    ctx,
			conf:   conf,
			conn:   c,
			roomID: boardID,
		}
		ws.ReadLoop(socket, service)
	}).ServeHTTP(ctx.Writer, ctx.Request)
}

func (ws *WebSocketServer) ReadLoop(socket *socket, service wsService) {
	buf := make([]byte, 1024)
	for {
		n, err := socket.conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				delete(ws.conns[socket.roomID], socket.conn)
				log.Info("connection closed: ", socket.conn)
				break
			}
			log.Error("websocket err: ", err)
			continue
		}

		msg := buf[:n]
		var message webSocketMessage
		json.Unmarshal(msg, &message)
		ws.HandleEvent(socket, service, message)
	}
}
