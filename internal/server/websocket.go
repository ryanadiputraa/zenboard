package server

import (
	"fmt"
	"io"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"

	log "github.com/sirupsen/logrus"
)

type WebSocketServer struct {
	conns map[string]map[*websocket.Conn]bool
	sync.Mutex
}

func (ws *WebSocketServer) HandleConnection(ctx *gin.Context) {
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

		ws.ReadLoop(c, boardID)
	}).ServeHTTP(ctx.Writer, ctx.Request)
}

func (ws *WebSocketServer) ReadLoop(conn *websocket.Conn, roomId string) {
	buf := make([]byte, 1024)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				delete(ws.conns[roomId], conn)
				log.Info("connection closed: ", conn)
				break
			}
			log.Error("websocket err: ", err)
			continue
		}
	}
}
