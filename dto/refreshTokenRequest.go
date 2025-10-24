package dto

import (
	"errors"

	"github.com/ashishjuyal/banking-auth/domain"
	"github.com/golang-jwt/jwt"
)

type RefreshTokenRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r RefreshTokenRequest) IsAccessTokenValid() *jwt.ValidationError {
	// 1. Token is invalid.
	// 2. Token is expired.
	_, err := jwt.Parse(r.AccessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		// invalid, expired errors differ only in message
		var vErr *jwt.ValidationError
		if errors.As(err, &vErr) {
			return vErr
		}
	}
	return nil
}
