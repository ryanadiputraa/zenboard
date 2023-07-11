package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryanadiputraa/zenboard/config"
	"github.com/ryanadiputraa/zenboard/internal/domain"
	"github.com/ryanadiputraa/zenboard/pkg/httpres"
	"github.com/ryanadiputraa/zenboard/pkg/jwt"
)

type boardController struct {
	conf    *config.Config
	service domain.BoardService
}

func NewBoardController(conf *config.Config, rg *gin.RouterGroup, service domain.BoardService) {
	c := boardController{
		conf:    conf,
		service: service,
	}

	rg.GET("/boards", c.GetUserBoards)
}

func (c *boardController) GetUserBoards(ctx *gin.Context) {
	userID, err := jwt.ExtractUserID(ctx, c.conf.JWT)
	if err != nil {
		httpres.HTTPErrorResponse(ctx, err)
		return
	}

	boards, err := c.service.GetUserBoards(ctx, userID)
	if err != nil {
		httpres.HTTPErrorResponse(ctx, err)
		return
	}

	httpres.HTTPSuccesResponse(ctx, http.StatusOK, boards)
}
