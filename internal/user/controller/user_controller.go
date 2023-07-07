package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryanadiputraa/zenboard/config"
	"github.com/ryanadiputraa/zenboard/internal/domain"
	"github.com/ryanadiputraa/zenboard/pkg/httpres"
	"github.com/ryanadiputraa/zenboard/pkg/jwt"
)

type userController struct {
	conf    *config.Config
	service domain.UserService
}

type findByIDReq struct {
	ID string `uri:"id" binding:"required"`
}

func NewUserController(conf *config.Config, rg *gin.RouterGroup, service domain.UserService) {
	c := userController{conf: conf, service: service}
	r := rg.Group("/users")

	r.GET("/", c.GetUserInfo)
	r.GET("/:id", c.FindUserByID)
}

func (c *userController) GetUserInfo(ctx *gin.Context) {
	userID, err := jwt.ExtractUserID(ctx, c.conf.JWT)
	if err != nil {
		httpres.HTTPErrorResponse(ctx, err)
		return
	}

	user, err := c.service.FindUserByID(ctx, userID)
	if err != nil {
		httpres.HTTPErrorResponse(ctx, err)
		return
	}

	httpres.HTTPSuccesResponse(ctx, http.StatusOK, user)
}

func (c *userController) FindUserByID(ctx *gin.Context) {
	var req findByIDReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		httpres.HTTPErrorResponse(ctx, err)
		return
	}

	user, err := c.service.FindUserByID(ctx, req.ID)
	if err != nil {
		httpres.HTTPErrorResponse(ctx, err)
		return
	}

	httpres.HTTPSuccesResponse(ctx, http.StatusOK, user)
}
