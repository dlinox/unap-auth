// adapter/http/response.go
package http

import (
	"github.com/gofiber/fiber/v2"
)

// Response representa una estructura genérica de respuesta
type Response struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// JSONResponse envía una respuesta JSON estandarizada con Fiber
func JSONResponse(c *fiber.Ctx, statusCode int, status string, data interface{}, message string, errors interface{}) error {
	response := Response{
		Status:  status,
		Data:    data,
		Message: message,
		Errors:  errors,
	}

	return c.Status(statusCode).JSON(response)
}
