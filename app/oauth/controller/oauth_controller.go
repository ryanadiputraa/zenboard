package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryanadiputraa/zenboard/domain"
	"github.com/ryanadiputraa/zenboard/pkg/httpres"
	"github.com/ryanadiputraa/zenboard/pkg/jwt"
	"github.com/ryanadiputraa/zenboard/pkg/oauth"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type oauthController struct {
	service      domain.OauthService
	userService  domain.UserService
	boardService domain.BoardService
}

type callbackReq struct {
	State string `form:"state" binding:"required"`
	Code  string `form:"code" binding:"required"`
}

func NewOauthController(
	rg *gin.RouterGroup,
	service domain.OauthService,
	userService domain.UserService,
	boardService domain.BoardService,
) {
	c := oauthController{
		service:      service,
		userService:  userService,
		boardService: boardService,
	}

	rg.GET("/login/google", c.LoginGoogle)
	rg.GET("/callback", c.Callback)
	rg.POST("/refresh", c.RefreshToken)
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
		log.Error("fail to bind oauth redirect url: ", err)
		return
	}

	userInfo, err := c.service.Callback(ctx, req.State, req.Code)
	if err != nil {
		oauth.RedirectWithError(ctx, err.Error())
		return
	}

	user, err := c.userService.CreateOrUpdateUserIfExists(ctx, domain.User{
		ID:            userInfo.ID,
		FirstName:     userInfo.FirstName,
		LastName:      userInfo.LastName,
		Email:         userInfo.Email,
		Picture:       userInfo.Picture,
		Locale:        userInfo.Locale,
		VerifiedEmail: userInfo.VerifiedEmail,
	})
	if err != nil {
		oauth.RedirectWithError(ctx, err.Error())
		return
	}

	c.boardService.InitBoard(ctx, user.ID)

	tokens, err := c.service.Login(ctx, user.ID)
	if err != nil {
		oauth.RedirectWithError(ctx, err.Error())
		return
	}
	oauth.RedirectWithJWTTokens(ctx, tokens)
}

func (c *oauthController) RefreshToken(ctx *gin.Context) {
	refreshToken, err := jwt.ExtractTokenFromAuthorizationHeader(ctx)
	if err != nil {
		httpres.HTTPErrorResponse(ctx, err)
		return
	}

	token, err := c.service.RefreshAccessToken(ctx, refreshToken)
	if err != nil {
		httpres.HTTPErrorResponse(ctx, err)
		return
	}

	httpres.HTTPSuccesResponse(ctx, http.StatusOK, token)
}
