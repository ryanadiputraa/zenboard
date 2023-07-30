package server

import (
	"encoding/json"

	"github.com/ryanadiputraa/zenboard/pkg/jwt"
	"golang.org/x/net/websocket"
)

type socketResponse struct {
	Key       string `json:"key"`
	IsSuccess bool   `json:"is_success"`
	Message   string `json:"message"`
	Data      any    `json:"data"`
}

type authPayload struct {
	AccessToken string `json:"access_token"`
}

type deleteTaskPayload struct {
	TaskID string `json:"task_id"`
}

func convertMsgData[T any](data any) (target T) {
	v, _ := json.Marshal(data)
	json.Unmarshal([]byte(v), &target)
	return
}

// TODO: refactor send message to broadcast to room id instead of single connection
func sendMessage(c *websocket.Conn, key string, isSuccess bool, message string, data any) {
	resp := socketResponse{
		Key:       key,
		IsSuccess: isSuccess,
		Message:   message,
		Data:      data,
	}
	msg, _ := json.Marshal(resp)
	c.Write(msg)
}

func (ws *WebSocketServer) HandleEvent(socket *socket, service wsService, msg webSocketEventMessage) {
	switch msg.Key {
	case "auth":
		data := convertMsgData[authPayload](msg.Data)
		ws.conns[socket.roomID][socket.conn] = data.AccessToken
		sendMessage(socket.conn, msg.Key, true, "user authenticated", nil)

	case "delete_task":
		token := ws.conns[socket.roomID][socket.conn]
		userID, err := jwt.ExtractUserIDFromJWTToken(socket.conf, token)
		if err != nil {
			sendMessage(socket.conn, msg.Key, false, err.Error(), nil)
			return
		}

		isAuthorized, err := service.boardService.CheckIsUserAuthorized(socket.ctx, socket.roomID, userID)
		if err != nil || !isAuthorized {
			sendMessage(socket.conn, msg.Key, false, err.Error(), nil)
			return
		}

		data := convertMsgData[deleteTaskPayload](msg.Data)
		if data.TaskID == "" {
			sendMessage(socket.conn, msg.Key, false, "invalid param", nil)
			return
		}

		err = service.taskService.DeleteTask(socket.ctx, data.TaskID)
		if err != nil {
			sendMessage(socket.conn, msg.Key, false, err.Error(), nil)
			return
		}
		sendMessage(socket.conn, msg.Key, true, "task deleted", data.TaskID)
	}
}
