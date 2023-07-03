package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/ryanadiputraa/zenboard/domain"
	"github.com/ryanadiputraa/zenboard/pkg/jwt"
	"github.com/ryanadiputraa/zenboard/pkg/oauth"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type oauthService struct{}

func NewOauthService() domain.OauthService {
	return &oauthService{}
}

func (s *oauthService) Callback(ctx context.Context, state, code string) (userInfo domain.UserInfo, err error) {
	if state != viper.GetString("OAUTH_STATE") {
		log.Warn("invalid oauth state: ", state)
		return userInfo, errors.New("invalid oauth state")
	}

	token, err := oauth.OauthConfig().Exchange(context.Background(), code)
	if err != nil {
		log.Error("fail to exchange oauth code with token: ", err)
		return userInfo, errors.New("fail to retrieve user token")
	}

	resp, err := http.Get(oauth.UserInfoURL + "?access_token=" + token.AccessToken)
	if err != nil {
		log.Error("fail to retrieve user info: ", err)
		return userInfo, errors.New("fail to retrieve user info")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("fail to read user info response body: ", err)
		return userInfo, errors.New("fail to retrieve user info")
	}
	log.Trace("oauth login: ", string(body))

	json.Unmarshal(body, &userInfo)
	if err != nil {
		log.Error("fail to parse user info response: ", err)
		return userInfo, errors.New("fail to retrieve user info")
	}

	return
}

func (s *oauthService) Login(ctx context.Context, userID any) (tokens domain.Tokens, err error) {
	tokens, err = jwt.GenerateAccessToken(userID)
	if err != nil {
		log.Error("fail to sign in user: ", err)
	}
	return
}

func (s *oauthService) RefreshAccessToken(ctx context.Context, refreshToken string) (tokens domain.Tokens, err error) {
	claims, err := jwt.ParseJWTClaims(refreshToken, true)
	if err != nil {
		log.Warn("fail to parse refresh token: ", err)
		return
	}

	tokens, err = jwt.GenerateAccessToken(claims.Sub)
	if err != nil {
		log.Error("fail refresh access token: ", err)
	}
	return
}
