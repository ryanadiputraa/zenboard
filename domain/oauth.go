package domain

import "context"

type OauthService interface {
	Callback(ctx context.Context, state, code string) (userInfo UserInfo, err error)
}

type UserInfo struct {
	ID            string `json:"id"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}
