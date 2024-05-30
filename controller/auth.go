package controller

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/po3nx/fgtest/database"
	"github.com/po3nx/fgtest/models"
	"golang.org/x/crypto/bcrypt"
)

// Login get user and password
func Login(c *fiber.Ctx) error {
	user := []models.User{}
	type LoginInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	username := input.Username
	pass := []byte(input.Password)
	//hashedPassword, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	//if err != nil {
    //    panic(err)
    //}
	if username == "" || pass == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	//db.Where(&User{Name: "user", Gender: "Male"}).First(&user)
	r :=database.DBConn.Model(&models.User{}).Where("username = ?", username).First(&user)
	if (r.RowsAffected == 0){
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	err := bcrypt.CompareHashAndPassword(user.Password, pass)
	if err != nil {
        return c.SendStatus(fiber.StatusUnauthorized)
    }
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": t})
}