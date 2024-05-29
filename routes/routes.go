package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/po3nx/fgtest/controller"
)

func SetUpRoutes(app *fiber.App) {
	app.Get("/hello", controller.Hello)
	app.Get("/allbooks", controller.AllBooks)
	app.Get("/book/:id", controller.GetBook)
	app.Post("/book", controller.AddBook)
	app.Put("/book/:id", controller.Update)
	app.Delete("/book/:id", controller.Delete)
}