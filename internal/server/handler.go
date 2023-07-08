package server

import (
	_userController "github.com/ryanadiputraa/zenboard/internal/user/controller"
	_userRepository "github.com/ryanadiputraa/zenboard/internal/user/repository"
	_userService "github.com/ryanadiputraa/zenboard/internal/user/service"

	_boardRepository "github.com/ryanadiputraa/zenboard/internal/board/repository"
	_boardService "github.com/ryanadiputraa/zenboard/internal/board/service"

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

	// oauth
	oauthSerivce := _oauthService.NewOauthService(s.conf)
	_oauthController.NewOauthController(s.conf, oauth, oauthSerivce, userService, boardService)
}
