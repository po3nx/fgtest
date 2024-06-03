package controller

import (
    "fmt"
	"time"
    "github.com/go-ldap/ldap"
	"github.com/po3nx/fgtest/config"
	"github.com/po3nx/fgtest/database"
	"github.com/po3nx/fgtest/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type LDAPConfig struct {
    Server   string
    Port     int
    BindDN   string
    Password string
    SearchDN string
}

func LoadLDAPConfig() *LDAPConfig {
	port, _ := strconv.Atoi(config.Config("LDAP_PORT"))
    return &LDAPConfig{
        Server:   config.Config("LDAP_SERVER"),
        Port:     port,
        BindDN:   config.Config("LDAP_BINDN"),
        Password: config.Config("LDAP_PASSWORD"),
        SearchDN: config.Config("LDAP_DN"),
    }
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLDAPData struct {
    ID       string
    Email    string
    Name     string
    FullName string
}

func AuthUsingLDAP(username, password string) ( *UserLDAPData, error) {
	ldapConfig := LoadLDAPConfig()
    l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapConfig.Server, ldapConfig.Port))
    if err != nil {
        return  nil, err
    }
    defer l.Close()

    err = l.Bind(ldapConfig.BindDN, ldapConfig.Password)
	if err != nil {
		return  nil, err
	}
	searchRequest := ldap.NewSearchRequest(
		ldapConfig.SearchDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(samaccountname=%s)", username),
		[]string{"dn", "cn", "samaccountname", "mail","telephonenumber","sn","givenname","distinguishedname","displayname"},
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		return  nil, err
	}

	if len(sr.Entries) == 0 {
		return  nil, fmt.Errorf("User not found")
	}
	entry := sr.Entries[0]

	err = l.Bind(entry.DN, password)
	if err != nil {
		return  nil, err
	}
	data := new(UserLDAPData)
	data.ID = username

	for _, attr := range entry.Attributes {
		switch attr.Name {
		case "sn":
			data.Name = attr.Values[0]
		case "mail":
			data.Email = attr.Values[0]
		case "cn":
			data.FullName = attr.Values[0]
		}
	}

	return data, nil
}
func LoginLDAP(c *fiber.Ctx) error {
	var loginReq LoginRequest

	if err := c.BodyParser(&loginReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request",
			"data":    err.Error(),
		})
	}

	data, err := AuthUsingLDAP(loginReq.Username, loginReq.Password)
	if err != nil {
		var user models.User
		r := database.DBConn.Model(&models.User{}).Where("username = ?", loginReq.Username).First(&user)
		if r.RowsAffected == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid username or password",
				"data":    nil,
			})
		}
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid username or password",
				"data":    nil,
			})
		}
		data = &UserLDAPData{
			ID:    user.Username,
			Email: user.Email,
		}
	}

	var user models.User
	r :=database.DBConn.Model(&models.User{}).Where("username = ?", data.ID).First(&user)
	pass := []byte(loginReq.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if r.Error != nil {
        // User not found, create a new user
        user = models.User{
            Username: data.ID,
            Email:    data.Email,
			Password :string(hashedPassword),
        }
        database.DBConn.Create(&user)
    } else {
        // User found, update the user data
        user.Email = data.Email
		user.Password = string(hashedPassword)
        database.DBConn.Save(&user)
    }
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = data.ID
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Success login",
		"data":    t,
	})
}