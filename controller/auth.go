package controller

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/po3nx/fgtest/database"
	"github.com/po3nx/fgtest/models"
	"golang.org/x/crypto/bcrypt"
	"github.com/po3nx/fgtest/utils"
)

// Login get user and password
func Login(c *fiber.Ctx) error {
	user := models.User{}
	type LoginInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	username := input.Username
	pass := input.Password
	
	//db.Where(&User{Name: "user", Gender: "Male"}).First(&user)
	r :=database.DBConn.Model(&models.User{}).Where("username = ?", username).First(&user)
	if (r.RowsAffected == 0){
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
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
	return utils.JSONResponse(c, "success", "Success login", t)
}

func Register(c *fiber.Ctx) error {
	type RegisterInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email string `json:"email"`
	}
	var input RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	username := input.Username
	pass := []byte(input.Password)
	email := input.Email
	hashedPassword, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
        panic(err)
    }
	user := models.User{Username :  username , Password : string(hashedPassword), Email : email}

	result := database.DBConn.Create(&user)
    if result.Error != nil {
		return utils.JSONResponse(c, "error", "Registration Failed", result.Error)
    }

	return utils.JSONResponse(c, "success", "User registered successfully", user)
}