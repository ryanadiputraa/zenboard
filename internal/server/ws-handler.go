package server

import (
	"encoding/json"

	"github.com/ryanadiputraa/zenboard/pkg/jwt"
)

type deleteTaskPayload struct {
	TaskID string `json:"task_id"`
}

func convertMsgData[T any](data any) (target T) {
	v, _ := json.Marshal(data)
	json.Unmarshal([]byte(v), &target)
	return
}

func (ws *WebSocketServer) HandleEvent(socket *socket, service wsService, msg webSocketEventMessage) {
	switch msg.Key {
	case "delete_task":
		userID, err := jwt.ExtractUserID(socket.ctx, socket.conf)
		if err != nil {
			socket.conn.Write([]byte(err.Error()))
			return
		}
		isAuthorized, err := service.boardService.CheckIsUserAuthorized(socket.ctx, socket.roomID, userID)
		if err != nil || !isAuthorized {
			socket.conn.Write([]byte(err.Error()))
			return
		}

		data := convertMsgData[deleteTaskPayload](msg.Data)
		if data.TaskID == "" {
			socket.conn.Write([]byte("invalid param"))
			return
		}

		err = service.taskService.DeleteTask(socket.ctx, data.TaskID)
		if err != nil {
			socket.conn.Write([]byte(err.Error()))
			return
		}
		socket.conn.Write([]byte("task deleted"))
	}
}
