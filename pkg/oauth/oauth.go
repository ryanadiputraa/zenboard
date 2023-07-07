package oauth

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ryanadiputraa/zenboard/config"
	"github.com/ryanadiputraa/zenboard/internal/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const UserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"

func OauthConfig(conf config.Oauth) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  conf.CallbackURL,
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile", "openid"},
		Endpoint:     google.Endpoint,
	}
}

func RedirectWithError(conf config.Oauth, ctx *gin.Context, err string) {
	redirectURL := conf.RedirectURL
	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?err=%s", redirectURL, err))
}

func RedirectWithJWTTokens(conf config.Oauth, ctx *gin.Context, tokens domain.JWTToken) {
	baseEedirectURL := conf.RedirectURL
	exp := strconv.FormatInt(tokens.ExpiresIn, 10)

	redirectURL := fmt.Sprintf(
		"%s?access_token=%s&expires_in=%s&refresh_token=%s",
		baseEedirectURL,
		tokens.AccessToken,
		exp,
		tokens.RefreshToken,
	)
	ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
