package oauth

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ryanadiputraa/zenboard/domain"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const UserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"

func OauthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  viper.GetString("OAUTH_CALLBACK_URL"),
		ClientID:     viper.GetString("OAUTH_CLIENT_ID"),
		ClientSecret: viper.GetString("OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile", "openid"},
		Endpoint:     google.Endpoint,
	}
}

func RedirectWithError(ctx *gin.Context, err string) {
	redirectURL := viper.GetString("OAUTH_REDIRECT_URL")
	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?err=%s", redirectURL, err))
}

func RedirectWithJWTTokens(ctx *gin.Context, tokens domain.Tokens) {
	baseEedirectURL := viper.GetString("OAUTH_REDIRECT_URL")
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
