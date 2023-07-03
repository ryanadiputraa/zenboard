package oauth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
