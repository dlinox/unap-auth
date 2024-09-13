package http

import (
	"unap-auth/usecase"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	AuthUsecase *usecase.AuthUsecase
}

// NewAuthHandler crea una nueva instancia de AuthHandler
func NewAuthHandler(authUsecase *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		AuthUsecase: authUsecase,
	}
}

func (h *AuthHandler) SignIn(c *fiber.Ctx) error {

	var req struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return JSONResponse(c, fiber.StatusBadRequest, "error", nil, "Invalid request body", nil)
	}

	if req.UserName == "" || req.Password == "" {
		return JSONResponse(c, fiber.StatusUnauthorized, "error", nil, "Invalid request body", nil)
	}

	token, err := h.AuthUsecase.Authenticate(req.UserName, req.Password)
	if err != nil {
		return JSONResponse(c, fiber.StatusUnauthorized, "error", nil, "email or password is incorrect", nil)
	}

	return JSONResponse(c, fiber.StatusOK, "success", token, "User authenticated successfully", nil)
}

func (h *AuthHandler) GetRoles(c *fiber.Ctx) error {
	// Extraer el user_id del contexto
	userAccountId, ok := c.Locals("uaid").(string)
	if !ok {
		return JSONResponse(c, fiber.StatusUnauthorized, "error", nil, "Invalid token payload", nil)
	}

	// Llamar al caso de uso para obtener los roles del usuario
	roles, err := h.AuthUsecase.GetRoles(userAccountId)
	if err != nil {

		return JSONResponse(c, fiber.StatusInternalServerError, "error", nil, "Could not retrieve roles", nil)
	}

	return JSONResponse(c, fiber.StatusOK, "success", roles, "Roles retrieved successfully", nil)
}

func (h *AuthHandler) GetModules(c *fiber.Ctx) error {

	// Extraer el user_id del contexto

	RoleId := c.Params("role_id")

	// Llamar al caso de uso para obtener los roles del usuario
	modules, err := h.AuthUsecase.GetModulesByRole(RoleId)

	if err != nil {
		return JSONResponse(c, fiber.StatusInternalServerError, "error", nil, "Could not retrieve modules", err.Error())
	}

	return JSONResponse(c, fiber.StatusOK, "success", modules, "Modules retrieved successfully", nil)
}

// const response = await http(token).post(`/auth/authorize`, data);
func (h *AuthHandler) AuthorizeToken(c *fiber.Ctx) error {

	var req struct {
		RoleId   string `json:"roleId"`
		ModuleId string `json:"moduleId"`
	}

	if err := c.BodyParser(&req); err != nil {
		return JSONResponse(c, fiber.StatusBadRequest, "error", nil, "Invalid request body", nil)
	}

	// Extraer el user_id del contexto
	userAccountId, ok := c.Locals("uaid").(string)

	if !ok {
		return JSONResponse(c, fiber.StatusUnauthorized, "error", nil, "Invalid token payload", nil)
	}

	// Llamar al caso de uso para obtener los roles del usuario
	token, err := h.AuthUsecase.AuthorizeToken(userAccountId, req.RoleId, req.ModuleId)
	if err != nil {
		return JSONResponse(c, fiber.StatusInternalServerError, "error", nil, "Could not authorize token", err.Error())
	}
	return JSONResponse(c, fiber.StatusOK, "success", token, "Token authorized successfully", nil)
}

// validateToken
func (h *AuthHandler) ValidateToken(c *fiber.Ctx) error {

	authorization := c.Get("Authorization")

	if authorization == "" {
		return JSONResponse(c, fiber.StatusUnauthorized, "error", nil, "No token provided", nil)
	}

	token := authorization[7:]

	valid := h.AuthUsecase.ValidateToken(token)

	if !valid {
		return JSONResponse(c, fiber.StatusUnauthorized, "error", token, "Invalid token", nil)
	}

	return JSONResponse(c, fiber.StatusOK, "success", nil, "Token is valid", nil)
}

func (h *AuthHandler) AuthMiddleware(c *fiber.Ctx) error {

	authorization := c.Get("Authorization")

	if authorization == "" {
		return JSONResponse(c, fiber.StatusUnauthorized, "error", nil, "No token provided", nil)
	}

	token := authorization[7:]

	if token == "" {
		return JSONResponse(c, fiber.StatusUnauthorized, "error", nil, "No token provided", nil)
	}

	auth, err := h.AuthUsecase.AuthMiddleware(token)

	if err != nil {
		return JSONResponse(c, fiber.StatusUnauthorized, "error", nil, "Invalid token", nil)
	}

	return JSONResponse(c, fiber.StatusOK, "success", auth, "Token is valid", nil)
}
