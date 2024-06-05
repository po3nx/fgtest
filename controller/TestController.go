package controller

import (
	"github.com/po3nx/fgtest/database"
	"github.com/po3nx/fgtest/models"
	"github.com/po3nx/fgtest/utils"
	"strconv"
	"github.com/gofiber/fiber/v2"
)

//Add
func Add(c *fiber.Ctx) error {
	book := new(models.Book)
	if err := c.BodyParser(book); err != nil {
		return utils.JSONResponse(c, "error", err.Error(), nil)
	}
	
	result := database.DBConn.Create(&book)
    if result.Error != nil {
        return utils.JSONResponse(c, "error", result.Error.Error(), nil)
    }

	return utils.JSONResponse(c, "success", "Book added successfully", book)
}

//GetByID
func GetByID(c *fiber.Ctx) error {
	book := []models.Book{}

	if err := database.DBConn.First(&book, c.Params("id")).Error; err != nil {
		return utils.JSONResponse(c, "error", err.Error(), nil)
	}

	return utils.JSONResponse(c, "success", "Book retrieved successfully", book)
}

//GetAll
func GetAll(c *fiber.Ctx) error {
	book := []models.Book{}

	if err := database.DBConn.Find(&book).Error; err != nil {
		return utils.JSONResponse(c, "error", err.Error(), nil)
	}

	return utils.JSONResponse(c, "success", "Books retrieved successfully", book)
}

//Update
func Update(c *fiber.Ctx) error {
	book := new(models.Book)
	if err := c.BodyParser(book); err != nil {
		return utils.JSONResponse(c, "error", err.Error(), nil)
	}
	id, _ := strconv.Atoi(c.Params("id"))

	result := database.DBConn.Model(&models.Book{}).Where("id = ?", id).Update("Title", book.Title)
    if result.Error != nil {
        return utils.JSONResponse(c, "error", result.Error.Error(), nil)
    }

	return utils.JSONResponse(c, "success", "Book updated successfully", nil)
}

//Delete
func Delete(c *fiber.Ctx) error {
	book := new(models.Book)

	id, _ := strconv.Atoi(c.Params("id"))

	if err := database.DBConn.Where("id = ?", id).Delete(&book).Error; err != nil {
		return utils.JSONResponse(c, "error", err.Error(), nil)
	}

	return utils.JSONResponse(c, "success", "Book deleted successfully", nil)
}