package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ryanadiputraa/zenboard/domain"
	"github.com/spf13/viper"
)

func GenerateAccessToken(userID any) (tokens domain.Tokens, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(viper.GetDuration("JWT_EXPIRES_IN")).Unix(),
	})
	tokenString, err := token.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		return
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(viper.GetDuration("JWT_REFRESH_EXPIRES_IN")).Unix(),
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(viper.GetString("JWT_REFRESH_SECRET")))
	if err != nil {
		return
	}

	tokens = domain.Tokens{
		AccessToken:  tokenString,
		ExpiresIn:    time.Now().Add(viper.GetDuration("JWT_EXPIRES_IN")).Unix(),
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

func ParseJWTClaims(tokenString string, isRefresh bool) (claims domain.JWTClaims, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return claims, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		var secret string
		if !isRefresh {
			secret = viper.GetString("JWT_SECRET")
		} else {
			secret = viper.GetString("JWT_REFRESH_SECRET")
		}
		return []byte(secret), nil
	})

	jwtClaim, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return claims, err
	}

	exp := jwtClaim["exp"].(float64)
	claims = domain.JWTClaims{
		Sub: jwtClaim["sub"].(string),
		Exp: exp,
	}

	return claims, nil
}

func ExtractUserID(ctx *gin.Context) (string, error) {
	token, err := ExtractTokenFromAuthorizationHeader(ctx)
	if err != nil {
		return "", err
	}

	claim, err := ParseJWTClaims(token, false)
	return claim.Sub, err
}
