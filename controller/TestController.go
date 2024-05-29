package controller

import (
	"github.com/po3nx/fgtest/database"
	"github.com/po3nx/fgtest/models"
	"strconv"
	"github.com/gofiber/fiber/v2"
)

//Hello
func Hello(c *fiber.Ctx) error {
	return c.SendString("fiber")
}

//Add
func Add(c *fiber.Ctx) error {
	test := new(models.Test)
	if err := c.BodyParser(test); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	
	result := database.DBConn.Create(&test)
    if result.Error != nil {
        return c.Status(400).JSON(result.Error) 
    }

	return c.Status(200).JSON(test)
}

//GetByID
func GetByID(c *fiber.Ctx) error {
	test := []models.Test{}

	database.DBConn.First(&test, c.Params("id"))

	return c.Status(200).JSON(test)
}

//GetAll
func GetAll(c *fiber.Ctx) error {
	test := []models.Test{}

	database.DBConn.Find(&test)

	return c.Status(200).JSON(test)
}

//Update
func Update(c *fiber.Ctx) error {
	test := new(models.Test)
	if err := c.BodyParser(test); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	id, _ := strconv.Atoi(c.Params("id"))

	database.DBConn.Model(&models.Test{}).Where("id = ?", id).Update("Username", test.Username)

	return c.Status(400).JSON("updated")
}

//Delete
func Delete(c *fiber.Ctx) error {
	test := new(models.Test)

	id, _ := strconv.Atoi(c.Params("id"))

	database.DBConn.Where("id = ?", id).Delete(&test)

	return c.Status(200).JSON("deleted")
}