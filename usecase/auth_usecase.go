package usecase

import (
	"time"
	"unap-auth/domain/repository"

	"github.com/golang-jwt/jwt/v5"
)

type AuthUsecase struct {
	UserAccountRepo repository.UserAccountRepository
	JwtSecret       string
}

func NewAuthUsecase(userAccountRepo repository.UserAccountRepository, jwtSecret string) *AuthUsecase {
	return &AuthUsecase{
		UserAccountRepo: userAccountRepo,
		JwtSecret:       jwtSecret,
	}
}

func (a *AuthUsecase) Authenticate(userName string, password string) (string, error) {
	user, err := a.UserAccountRepo.FindAndValidateUserAccount(userName, password)
	if err != nil || user == nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uaid": user.Id,
		"exp":  time.Now().Add(time.Hour * 2).Unix(),
	})

	signedToken, err := token.SignedString([]byte(a.JwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
