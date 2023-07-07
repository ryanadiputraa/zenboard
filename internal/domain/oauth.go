package domain

import "context"

type OauthService interface {
	Callback(ctx context.Context, state, code string) (userInfo UserInfo, err error)
	Login(ctx context.Context, userID any) (JWTToken, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (JWTToken, error)
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

type JWTToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type JWTClaims struct {
	Sub string  `json:"sub"`
	Exp float64 `json:"exp"`
}
