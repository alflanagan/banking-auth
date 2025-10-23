package service

import (
	"fmt"

	"github.com/alflanagan/banking-lib/errs"
	"github.com/alflanagan/banking-lib/logger"
	"github.com/ashishjuyal/banking-auth/domain"
	"github.com/ashishjuyal/banking-auth/dto"
	"github.com/golang-jwt/jwt"
)

type AuthService interface {
	Login(dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
	Verify(urlParams map[string]string) *errs.AppError
}

type DefaultAuthService struct {
	repo            domain.AuthRepository
	rolePermissions domain.RolePermissions
}

func (s DefaultAuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, *errs.AppError) {
	var appErr *errs.AppError
	var login *domain.Login
	var accessToken string

	if login, appErr = s.repo.FindBy(req.Username, req.Password); appErr != nil {
		return nil, appErr
	}
	claims := login.ClaimsForAccessToken()
	authToken := domain.NewAuthToken(claims)
	if accessToken, appErr = authToken.NewAccessToken(); appErr != nil {
		return nil, appErr
	}
	return &dto.LoginResponse{AccessToken: accessToken}, nil
}

func (s DefaultAuthService) Verify(urlParams map[string]string) *errs.AppError {
	// convert the string token to JWT struct
	if jwtToken, err := jwtTokenFromString(urlParams["token"]); err != nil {
		return errs.NewAuthenticationError(err.Error())
	} else {
		/*
		   Checking the validity of the token, this verifies the expiry
		   time and the signature of the token
		*/
		if jwtToken.Valid {
			// type cast the token claims to jwt.MapClaims
			claims := jwtToken.Claims.(*domain.Claims)
			/* if Role is user then check if the account_id and customer_id
			   coming in the URL belongs to the same token
			*/
			if claims.IsUserRole() {
				if !claims.IsRequestVerifiedWithTokenClaims(urlParams) {
					return errs.NewAuthorizationError("user role not verified")
				}
			}
			// verify of the role is authorized to use the route
			isAuthorized := s.rolePermissions.IsAuthorizedFor(claims.Role, urlParams["routeName"])
			if !isAuthorized {
				return errs.NewAuthorizationError(fmt.Sprintf("user not authorized for %s role", claims.Role))
			}
			return nil
		} else {
			return errs.NewAuthorizationError("invalid token")
		}
	}
}

func jwtTokenFromString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		logger.Error("Error while parsing token: " + err.Error())
		return nil, err
	}
	return token, nil
}

func NewLoginService(repo domain.AuthRepository, permissions domain.RolePermissions) DefaultAuthService {
	return DefaultAuthService{repo, permissions}
}
