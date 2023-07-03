package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/ryanadiputraa/zenboard/domain"
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
		log.Error("invalid oauth state: ", state)
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
		log.Error("fail to retrieve user info: ", err)
		return userInfo, errors.New("fail to retrieve user info")
	}
	log.Info("oauth login: ", string(body))

	json.Unmarshal(body, &userInfo)
	if err != nil {
		log.Error("fail to decode response body from user info: ", err)
		return userInfo, errors.New("fail to retrieve user info")
	}

	return
}
