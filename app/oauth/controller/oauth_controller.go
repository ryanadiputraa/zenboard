package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryanadiputraa/zenboard/domain"
	"github.com/ryanadiputraa/zenboard/pkg/oauth"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type oauthController struct {
	service domain.OauthService
}

type callbackReq struct {
	State string `form:"state" binding:"required"`
	Code  string `form:"code" binding:"required"`
}

func NewOauthController(rg *gin.RouterGroup, service domain.OauthService) {
	c := oauthController{service: service}

	rg.GET("/login/google", c.LoginGoogle)
	rg.GET("/callback", c.Callback)
}

func (c *oauthController) LoginGoogle(ctx *gin.Context) {
	state := viper.GetString("OAUTH_STATE")
	url := oauth.OauthConfig().AuthCodeURL(state, oauth2.SetAuthURLParam("prompt", "select_account"))
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (c *oauthController) Callback(ctx *gin.Context) {
	var req callbackReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		oauth.RedirectWithError(ctx, err.Error())
		return
	}

	userInfo, err := c.service.Callback(ctx, req.State, req.Code)
	if err != nil {
		oauth.RedirectWithError(ctx, err.Error())
		return
	}

	// TODO: save user info

	tokens, err := c.service.Login(ctx, userInfo.ID)
	if err != nil {
		oauth.RedirectWithError(ctx, err.Error())
		return
	}
	oauth.RedirectWithJWTTokens(ctx, tokens)
}
