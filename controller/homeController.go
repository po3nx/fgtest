package controller

import (
	"github.com/po3nx/fgtest/utils"
	"github.com/gofiber/fiber/v2"
)

//Home
func Home(c *fiber.Ctx) error {
	return utils.JSONResponse(c, "success", "This is default root API Endpoint", nil)
}
