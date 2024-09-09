package mysql

import (
	"database/sql"
	"errors"
	"time"
	"unap-auth/domain/model"

	"golang.org/x/crypto/bcrypt"
)

type MySQLUserAccountRepository struct {
	DB *sql.DB
}

var (
	ErrUserNotFound    = errors.New("user account not found")
	ErrInvalidPassword = errors.New("username or password is incorrect")
	ErrInactiveAccount = errors.New("user account is inactive")
	ErrorUnauthorized  = errors.New("Unauthorized")
)

func (repo *MySQLUserAccountRepository) FindAndValidateUserAccount(userName string, password string) (*model.UserAccount, error) {
	query := `
			SELECT Id, UserName, Password, StartDate, CaducityDate, Status 
			FROM Core_UserAccounts 
			WHERE UserName = ?
			`
	row := repo.DB.QueryRow(query, userName)

	user := &model.UserAccount{}

	var startDateStr, caducityDateStr string

	err := row.Scan(&user.Id, &user.Username, &user.Password, &startDateStr, &caducityDateStr, &user.Status)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrInvalidPassword
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, err
	}

	caducityDate, err := time.Parse("2006-01-02", caducityDateStr)
	if err != nil {
		return nil, err
	}

	if user.Status != 7 || startDate.After(time.Now()) || caducityDate.Before(time.Now()) {
		return nil, ErrInactiveAccount
	}

	return user, nil
}

func (repo *MySQLUserAccountRepository) GetRoles(userAcountId string) ([]model.Role, error) {
	query := `
		SELECT cr.Id, cr.Name, cr.Description
		FROM core_roles cr
		INNER JOIN core_behavior cb ON cb.RoleId = cr.Id
		WHERE cb.UserAccountId = ? AND cr.Status = 7;
	`

	// Ejecutar la consulta con el userID
	rows, err := repo.DB.Query(query, userAcountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []model.Role

	for rows.Next() {
		role := model.Role{}
		err := rows.Scan(&role.Id, &role.Name, &role.Description)
		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	return roles, nil

}

func (repo *MySQLUserAccountRepository) GetModulesByRole(roleId string) ([]model.Module, error) {

	query := `
		SELECT
			cm.Id,
			cm.Name, 
			cm.Route, 
			cm.Description,
			IfNull(cm.Icon, "") as Icon
		FROM core_modules cm
		JOIN core_resources cr ON cr.ResourceId = cm.Id
		JOIN core_roleaccesses cra ON cra.ResourceId = cr.Id
		WHERE cra.RoleId = ?
		AND cr.ResourceParentId IS NULL
		AND cm.Status = 7;`

	rows, err := repo.DB.Query(query, roleId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var modules []model.Module

	for rows.Next() {
		module := model.Module{}
		err := rows.Scan(&module.Id, &module.Name, &module.Route, &module.Description, &module.Icon)
		if err != nil {
			return nil, err
		}
		modules = append(modules, module)
	}

	return modules, nil
}

func (repo *MySQLUserAccountRepository) AuthorizeToken(userAccountId string, roleId string, moduleId string) (*model.AuthorizeTokenClaims, error) {
	query := `
		SELECT cr.Id
		FROM core_roles cr
		INNER JOIN core_behavior cb ON cb.RoleId = cr.Id
		WHERE cb.UserAccountId = ? AND cr.Id = ? AND cr.Status = 7 
	`

	//validar que el usuario tenga el rol, si no, retornar error con el mensaje "No autorizado"
	row := repo.DB.QueryRow(query, userAccountId, roleId)
	if err := row.Scan(&roleId); err != nil {
		return nil, ErrorUnauthorized
	}

	//validar que el modulo pertenezca al rol, si no, retornar error con el mensaje "No autorizado"
	query = `
		SELECT cm.Id
		FROM core_modules cm
		JOIN core_resources cr ON cr.ResourceId = cm.Id
		JOIN core_roleaccesses cra ON cra.ResourceId = cr.Id
		WHERE cra.RoleId = ? AND cm.Id = ? AND cr.ResourceParentId IS NULL AND cm.Status = 7;`

	row = repo.DB.QueryRow(query, roleId, moduleId)
	if err := row.Scan(&moduleId); err != nil {
		return nil, ErrorUnauthorized
	}

	//generar el token con el modulo y el rol
	tokenClaims := model.AuthorizeTokenClaims{
		UserAccountId: userAccountId,
		RoleId:        roleId,
		ModuleId:      moduleId,
	}

	return &tokenClaims, nil

}

func (repo *MySQLUserAccountRepository) ValidateToken(token string) (bool, error) {
	return true, nil
}
