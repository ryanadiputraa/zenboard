package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/ryanadiputraa/zenboard/config"
	"github.com/ryanadiputraa/zenboard/internal/domain"
	"github.com/ryanadiputraa/zenboard/pkg/jwt"
	"github.com/ryanadiputraa/zenboard/pkg/oauth"
	log "github.com/sirupsen/logrus"
)

type oauthService struct {
	conf *config.Config
}

func NewOauthService(conf *config.Config) domain.OauthService {
	return &oauthService{conf: conf}
}

func (s *oauthService) Callback(ctx context.Context, state, code string) (userInfo domain.UserInfo, err error) {
	if state != s.conf.Oauth.State {
		log.Warn("invalid oauth state: ", state)
		return userInfo, errors.New("invalid oauth state")
	}

	token, err := oauth.OauthConfig(s.conf.Oauth).Exchange(ctx, code)
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
	log.Info("oauth login: ", string(body))

	json.Unmarshal(body, &userInfo)
	if err != nil {
		log.Error("fail to parse user info response: ", err)
		return userInfo, errors.New("fail to retrieve user info")
	}

	return
}

func (s *oauthService) Login(ctx context.Context, userID any) (tokens domain.JWTToken, err error) {
	tokens, err = jwt.GenerateAccessToken(s.conf.JWT, userID)
	if err != nil {
		log.Error("fail to sign in user: ", err)
	}
	return
}

func (s *oauthService) RefreshAccessToken(ctx context.Context, refreshToken string) (tokens domain.JWTToken, err error) {
	claims, err := jwt.ParseJWTClaims(s.conf.JWT, refreshToken, true)
	if err != nil {
		log.Warn("fail to parse refresh token: ", err)
		return
	}

	tokens, err = jwt.GenerateAccessToken(s.conf.JWT, claims.Sub)
	if err != nil {
		log.Error("fail refresh access token: ", err)
	}
	return
}
