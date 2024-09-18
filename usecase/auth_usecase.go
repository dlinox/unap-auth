package usecase

import (
	"time"
	"unap-auth/domain/model"
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

func (a *AuthUsecase) GetRoles(userAccountId string) ([]model.Role, error) {
	return a.UserAccountRepo.GetRoles(userAccountId)
}

func (a *AuthUsecase) GetModulesByRole(roleId string) ([]model.Module, error) {
	return a.UserAccountRepo.GetModulesByRole(roleId)
}

func (a *AuthUsecase) AuthorizeToken(userAccountId string, roleId string, moduleId string) (string, error) {

	tokenClaims, err := a.UserAccountRepo.AuthorizeToken(userAccountId, roleId, moduleId)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uaid": tokenClaims.UserAccountId,
		"rid":  tokenClaims.RoleId,
		"mid":  tokenClaims.ModuleId,
		"bid":  tokenClaims.BehaviorId,
		"exp":  time.Now().Add(time.Hour * 2).Unix(),
	})

	signedToken, err := token.SignedString([]byte(a.JwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (a *AuthUsecase) ValidateToken(tokenString string) bool {

	// ParseToken parses and validates the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verificar que el método de firma del token es HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil
		}
		return []byte(a.JwtSecret), nil
	})

	if err != nil {
		return false
	}

	// Verificar si el token es válido
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true
	}

	return false
}

func (a *AuthUsecase) AuthMiddleware(tokenString string) (*model.Auth, error) {

	// ParseToken parses and validates the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verificar que el método de firma del token es HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil
		}
		return []byte(a.JwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Verificar si el token es válido
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		permissions, err := a.UserAccountRepo.AuthMiddleware(claims["rid"].(string))

		auth := &model.Auth{
			UserAccountId: claims["uaid"].(string),
			RoleId:        claims["rid"].(string),
			ModuleId:      claims["mid"].(string),
			BehaviorId:    claims["bid"].(string),
			Permissions:   permissions,
		}
		if err != nil {
			return nil, err
		}

		return auth, nil

	}
	return nil, nil
}
