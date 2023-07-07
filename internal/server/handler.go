package server

import (
	_userController "github.com/ryanadiputraa/zenboard/app/user/controller"
	_userRepository "github.com/ryanadiputraa/zenboard/app/user/repository"
	_userService "github.com/ryanadiputraa/zenboard/app/user/service"
	db "github.com/ryanadiputraa/zenboard/pkg/db/sqlc"

	_boardRepository "github.com/ryanadiputraa/zenboard/app/board/repository"
	_boardService "github.com/ryanadiputraa/zenboard/app/board/service"

	_oauthController "github.com/ryanadiputraa/zenboard/app/oauth/controller"
	_oauthService "github.com/ryanadiputraa/zenboard/app/oauth/service"
)

func (s *Server) MapHandlers() {
	api := s.gin.Group("/api")
	oauth := s.gin.Group("/oauth")

	DB := db.New(s.db)
	Tx := db.NewTx(s.db)

	// user
	userRepository := _userRepository.NewUserRepository(DB)
	userService := _userService.NewUserService(userRepository)
	_userController.NewUserController(s.conf, api, userService)

	// board
	boardRepository := _boardRepository.NewBoardRepository(DB, Tx)
	boardService := _boardService.NewBoardService(boardRepository)

	// oauth
	oauthSerivce := _oauthService.NewOauthService(s.conf)
	_oauthController.NewOauthController(s.conf, oauth, oauthSerivce, userService, boardService)
}
