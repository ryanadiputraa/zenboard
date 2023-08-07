package server

import (
	"encoding/json"
	"errors"

	"github.com/ryanadiputraa/zenboard/internal/domain"
	"github.com/ryanadiputraa/zenboard/pkg/jwt"
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

func validateUserAndMessageData[T any](ws *WebSocketServer, socket *socket, boardService domain.BoardService, msgKey string, data any) (target T) {
	err := ws.validateUser(socket, boardService)
	if err != nil {
		ws.SendMessage(socket, msgKey, false, err.Error(), nil)
		return
	}
	target = convertMsgData[T](data)
	return
}

func (ws *WebSocketServer) HandleEvent(socket *socket, service wsService, msg webSocketEventMessage) {
	switch msg.Key {
	case "auth":
		data := convertMsgData[authPayload](msg.Data)
		ws.conns[socket.roomID][socket.conn] = data.AccessToken
		ws.SendMessage(socket, msg.Key, true, "user authenticated", nil)

	case "change_project_name":
		data := validateUserAndMessageData[changeProjectNamePayload](ws, socket, service.boardService, msg.Key, msg.Data)
		if data.Name == "" {
			ws.SendMessage(socket, msg.Key, false, "invalid param", nil)
			return
		}
		board, err := service.boardService.ChangeProjectName(socket.ctx, socket.roomID, data.Name)
		if err != nil {
			ws.SendMessage(socket, msg.Key, false, err.Error(), nil)
			return
		}
		ws.Broadcast(socket, msg.Key, true, "project name changed", board)

	case "create_task":
		data := validateUserAndMessageData[createTaskPayload](ws, socket, service.boardService, msg.Key, msg.Data)
		if data.TaskName == "" {
			ws.SendMessage(socket, msg.Key, false, "invalid param", nil)
			return
		}
		task, err := service.taskService.AddBoardTask(socket.ctx, socket.roomID, data.TaskName)
		if err != nil {
			ws.SendMessage(socket, msg.Key, false, err.Error(), nil)
			return
		}
		ws.Broadcast(socket, msg.Key, true, "task created", task)

	case "delete_task":
		data := validateUserAndMessageData[deleteTaskPayload](ws, socket, service.boardService, msg.Key, msg.Data)
		if data.TaskID == "" {
			ws.SendMessage(socket, msg.Key, false, "invalid param", nil)
			return
		}
		err := service.taskService.DeleteTask(socket.ctx, data.TaskID)
		if err != nil {
			ws.SendMessage(socket, msg.Key, false, err.Error(), nil)
			return
		}
		ws.Broadcast(socket, msg.Key, true, "task deleted", data.TaskID)

	case "reorder_task":
		data := validateUserAndMessageData[[]domain.TaskReorderDTO](ws, socket, service.boardService, msg.Key, msg.Data)
		if len(data) == 0 {
			ws.SendMessage(socket, msg.Key, false, "invalid param", nil)
			return
		}
		tasks, err := service.taskService.UpdateOrder(socket.ctx, data)
		if err != nil {
			ws.SendMessage(socket, msg.Key, false, err.Error(), nil)
			return
		}
		ws.Broadcast(socket, msg.Key, true, "task reordered", tasks)
	}
}
