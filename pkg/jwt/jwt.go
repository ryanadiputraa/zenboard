package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ryanadiputraa/zenboard/config"
	"github.com/ryanadiputraa/zenboard/domain"
)

func GenerateAccessToken(conf config.JWT, userID any) (tokens domain.JWTToken, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(conf.ExpiresIn).Unix(),
	})
	tokenString, err := token.SignedString([]byte(conf.Secret))
	if err != nil {
		return
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(conf.RefreshExpiresIn).Unix(),
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(conf.RefreshSecret))
	if err != nil {
		return
	}

	tokens = domain.JWTToken{
		AccessToken:  tokenString,
		ExpiresIn:    time.Now().Add(conf.ExpiresIn).Unix(),
		RefreshToken: refreshTokenString,
	}

	return
}

func ExtractTokenFromAuthorizationHeader(ctx *gin.Context) (token string, err error) {
	t := ctx.GetHeader("Authorization")
	if len(t) == 0 {
		return "", errors.New("missing authorization header")
	}

	h := strings.Split(t, " ")
	if len(h) < 2 || h[0] != "Bearer" {
		return "", errors.New("invalid token format")
	}

	token = h[1]
	return
}

func ParseJWTClaims(conf config.JWT, tokenString string, isRefresh bool) (claims domain.JWTClaims, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		var secret string
		if !isRefresh {
			secret = conf.Secret
		} else {
			secret = conf.RefreshSecret
		}
		return []byte(secret), nil
	})

	if token == nil {
		return
	}

	jwtClaim, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return
	}

	exp := jwtClaim["exp"].(float64)
	claims = domain.JWTClaims{
		Sub: jwtClaim["sub"].(string),
		Exp: exp,
	}

	return
}

func ExtractUserID(ctx *gin.Context, conf config.JWT) (string, error) {
	token, err := ExtractTokenFromAuthorizationHeader(ctx)
	if err != nil {
		return "", err
	}

	claim, err := ParseJWTClaims(conf, token, false)
	return claim.Sub, err
}
