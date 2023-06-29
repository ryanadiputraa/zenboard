package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryanadiputraa/zenboard/domain"
	"github.com/ryanadiputraa/zenboard/pkg/httpres"
)

type userController struct {
	service domain.UserService
}

type createUserReq struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Picture   string `json:"picture"`
	Locale    string `json:"locale" binding:"required"`
}

type findByIDReq struct {
	ID string `uri:"id" binding:"required"`
}

func NewUserController(rg *gin.RouterGroup, service domain.UserService) {
	c := userController{service: service}
	r := rg.Group("/users")

	r.POST("/", c.CreateUser)
	r.GET("/:id", c.FindUserByID)
}

func (c *userController) CreateUser(ctx *gin.Context) {
	var req createUserReq
	if err := ctx.ShouldBind(&req); err != nil {
		httpres.HTTPErrorResponse(ctx, err)
		return
	}

	user, err := c.service.CreateOrUpdateUserIfExists(ctx, domain.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Picture:   req.Picture,
		Locale:    req.Locale,
	})
	if err != nil {
		httpres.HTTPErrorResponse(ctx, err)
		return
	}

	httpres.HTTPSuccesResponse(ctx, http.StatusCreated, user)
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
