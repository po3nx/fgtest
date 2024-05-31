package utils

import "github.com/gofiber/fiber/v2"

func JSONResponse(c *fiber.Ctx, status string, message string, data interface{}) error {
	return c.JSON(fiber.Map{
		"status":  status,
		"message": message,
		"data":    data,
	})
}