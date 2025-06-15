package utils

import "github.com/gofiber/fiber/v2"

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(Response{
		Status:  "error",
		Message: message,
	})
}
