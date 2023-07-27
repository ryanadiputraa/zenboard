package server

import (
	"github.com/gin-gonic/gin"
	_userController "github.com/ryanadiputraa/zenboard/internal/user/controller"
	_userRepository "github.com/ryanadiputraa/zenboard/internal/user/repository"
	_userService "github.com/ryanadiputraa/zenboard/internal/user/service"

	_boardController "github.com/ryanadiputraa/zenboard/internal/board/controller"
	_boardRepository "github.com/ryanadiputraa/zenboard/internal/board/repository"
	_boardService "github.com/ryanadiputraa/zenboard/internal/board/service"

	_taskController "github.com/ryanadiputraa/zenboard/internal/task/controller"
	_taskRepository "github.com/ryanadiputraa/zenboard/internal/task/repository"
	_taskService "github.com/ryanadiputraa/zenboard/internal/task/service"

	_oauthController "github.com/ryanadiputraa/zenboard/internal/oauth/controller"
	_oauthService "github.com/ryanadiputraa/zenboard/internal/oauth/service"
)

func (s *Server) MapHandlers() {
	api := s.gin.Group("/api")
	oauth := s.gin.Group("/oauth")

	// user
	userRepository := _userRepository.NewUserRepository(s.db)
	userService := _userService.NewUserService(userRepository)
	_userController.NewUserController(s.conf, api, userService)

	// board
	boardRepository := _boardRepository.NewBoardRepository(s.db)
	boardService := _boardService.NewBoardService(boardRepository)
	_boardController.NewBoardController(s.conf, api, boardService)

	// task
	taskRepository := _taskRepository.NewTaskRepository(s.db)
	taskService := _taskService.NewTaskService(taskRepository)
	_taskController.NewTaskController(s.conf, api, taskService, boardService)

	// oauth
	oauthSerivce := _oauthService.NewOauthService(s.conf)
	_oauthController.NewOauthController(s.conf, oauth, oauthSerivce, userService, boardService)

	wsService := wsService{
		boardService: boardService,
	}
	// websocket
	s.gin.GET("/ws", func(ctx *gin.Context) {
		s.ws.HandleConnection(ctx, s.conf.JWT, wsService)
	})
}
