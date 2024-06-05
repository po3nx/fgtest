package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/po3nx/fgtest/controller"
	"github.com/po3nx/fgtest/middleware"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetUpRoutes(app *fiber.App) {

	api := app.Group("/api", logger.New())
	api.Get("/", controller.Hello)

	auth := api.Group("/auth")
	auth.Post("/login", controller.LoginLDAP)
	auth.Post("/register", controller.Register)

	product := api.Group("/product")
	product.Get("/allbooks", controller.GetAll)
	product.Get("/book/:id", controller.GetByID)
	product.Post("/book", middleware.Protected(), controller.Add)
	product.Put("/book/:id", middleware.Protected(), controller.Update)
	product.Delete("/book/:id", middleware.Protected(), controller.Delete)
}