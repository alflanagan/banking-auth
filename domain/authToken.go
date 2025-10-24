package domain

import (
	"github.com/alflanagan/banking-lib/errs"
	"github.com/alflanagan/banking-lib/logger"
	"github.com/golang-jwt/jwt"
)

type AuthToken struct {
	token *jwt.Token
}

func (t AuthToken) NewAccessToken() (string, *errs.AppError) {
	signedString, err := t.token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		logger.Error("Failed while signing access token: " + err.Error())
		return "", errs.NewUnexpectedError("cannot generate access token")
	}
	return signedString, nil
}

func (t AuthToken) newRefreshToken() (string, *errs.AppError) {
	c := t.token.Claims.(Claims)
	refreshToken := c.RefreshTokenClaims()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshToken)
	signedString, err := token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		logger.Error("Failed while signing refresh token: " + err.Error())
		return "", errs.NewUnexpectedError("cannot generate refresh token")
	}
	return signedString, nil
}

func NewAuthToken(claims Claims) AuthToken {
	return AuthToken{
		token: jwt.NewWithClaims(jwt.SigningMethodHS256, claims),
	}
}
