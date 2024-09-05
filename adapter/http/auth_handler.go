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
