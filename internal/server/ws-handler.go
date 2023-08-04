package server

import (
	"encoding/json"
	"errors"

	"github.com/ryanadiputraa/zenboard/internal/domain"
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

type changeProjectNamePayload struct {
	Name string `json:"name"`
}

type createTaskPayload struct {
	BoardID  string `json:"board_id"`
	TaskName string `json:"task_name"`
}

type deleteTaskPayload struct {
	TaskID string `json:"task_id"`
}

func convertMsgData[T any](data any) (target T) {
	v, _ := json.Marshal(data)
	json.Unmarshal([]byte(v), &target)
	return
}

func (ws *WebSocketServer) broadcast(roomID string, c *websocket.Conn, key string, isSuccess bool, message string, data any) {
	resp := socketResponse{
		Key:       key,
		IsSuccess: isSuccess,
		Message:   message,
		Data:      data,
	}
	msg, _ := json.Marshal(resp)
	for conn := range ws.conns[roomID] {
		conn.Write(msg)
	}
}

func (ws *WebSocketServer) validateUser(socket *socket, boardService domain.BoardService) (err error) {
	token := ws.conns[socket.roomID][socket.conn]
	userID, err := jwt.ExtractUserIDFromJWTToken(socket.conf, token)
	if err != nil {
		return
	}
	isAuthorized, err := boardService.CheckIsUserAuthorized(socket.ctx, socket.roomID, userID)
	if !isAuthorized {
		err = errors.New("forbidden")
	}
	return
}

func (ws *WebSocketServer) HandleEvent(socket *socket, service wsService, msg webSocketEventMessage) {
	switch msg.Key {
	case "auth":
		data := convertMsgData[authPayload](msg.Data)
		ws.conns[socket.roomID][socket.conn] = data.AccessToken
		ws.broadcast(socket.roomID, socket.conn, msg.Key, true, "user authenticated", nil)

	case "change_project_name":
		err := ws.validateUser(socket, service.boardService)
		if err != nil {
			ws.broadcast(socket.roomID, socket.conn, msg.Key, false, err.Error(), nil)
			return
		}
		data := convertMsgData[changeProjectNamePayload](msg.Data)
		if data.Name == "" {
			ws.broadcast(socket.roomID, socket.conn, msg.Key, false, "invalid param", nil)
			return
		}
		board, err := service.boardService.ChangeProjectName(socket.ctx, socket.roomID, data.Name)
		if err != nil {
			ws.broadcast(socket.roomID, socket.conn, msg.Key, false, err.Error(), nil)
			return
		}
		ws.broadcast(socket.roomID, socket.conn, msg.Key, true, "project name changed", board)

	case "create_task":
		err := ws.validateUser(socket, service.boardService)
		if err != nil {
			ws.broadcast(socket.roomID, socket.conn, msg.Key, false, err.Error(), nil)
			return
		}
		data := convertMsgData[createTaskPayload](msg.Data)
		if data.BoardID == "" || data.TaskName == "" {
			ws.broadcast(socket.roomID, socket.conn, msg.Key, false, "invalid param", nil)
			return
		}

		task, err := service.taskService.AddBoardTask(socket.ctx, data.BoardID, data.TaskName)
		if err != nil {
			ws.broadcast(socket.roomID, socket.conn, msg.Key, false, err.Error(), nil)
			return
		}
		ws.broadcast(socket.roomID, socket.conn, msg.Key, true, "task created", task)

	case "delete_task":
		err := ws.validateUser(socket, service.boardService)
		if err != nil {
			ws.broadcast(socket.roomID, socket.conn, msg.Key, false, err.Error(), nil)
		}

		data := convertMsgData[deleteTaskPayload](msg.Data)
		if data.TaskID == "" {
			ws.broadcast(socket.roomID, socket.conn, msg.Key, false, "invalid param", nil)
			return
		}

		err = service.taskService.DeleteTask(socket.ctx, data.TaskID)
		if err != nil {
			ws.broadcast(socket.roomID, socket.conn, msg.Key, false, err.Error(), nil)
			return
		}
		ws.broadcast(socket.roomID, socket.conn, msg.Key, true, "task deleted", data.TaskID)
	}
}
