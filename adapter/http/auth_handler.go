package http

import (
	"time"
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

	startTime := time.Now()

	var req struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}

	// Analizar el cuerpo de la solicitud
	if err := c.BodyParser(&req); err != nil {
		// Enviar un mensaje de error más detallado
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if req.UserName == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	token, err := h.AuthUsecase.Authenticate(req.UserName, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	return c.JSON(fiber.Map{"token": token, "duration": duration.Seconds()})
}
