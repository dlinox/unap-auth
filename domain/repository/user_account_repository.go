package repository

import "unap-auth/domain/model"

type UserAccountRepository interface {
	FindAndValidateUserAccount(username string, password string) (*model.UserAccount, error)
	GetRoles(userAccountId string) ([]model.Role, error)
	GetModulesByRole(roleId string) ([]model.Module, error)
	AuthorizeToken(userAccountId string, roleId string, moduleId string) (*model.AuthorizeTokenClaims, error)
	ValidateToken(token string) (bool, error)
	AuthMiddleware(roleId string) (string, error)
}
