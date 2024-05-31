package controller

import (
    "fmt"
	"time"
    "github.com/go-ldap/ldap"
	"github.com/po3nx/fgtest/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

func AuthUsingLDAP(username, password string) (bool, *UserLDAPData, error) {
	ldapConfig := LoadLDAPConfig()
    l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapConfig.Server, ldapConfig.Port))
    if err != nil {
        return false, nil, err
    }
    defer l.Close()

    err = l.Bind(ldapConfig.BindDN, ldapConfig.Password)
	if err != nil {
		return false, nil, err
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
		return false, nil, err
	}

	if len(sr.Entries) == 0 {
		return false, nil, fmt.Errorf("User not found")
	}
	entry := sr.Entries[0]

	err = l.Bind(entry.DN, password)
	if err != nil {
		return false, nil, err
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

	return true, data, nil
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

	success, data, err := AuthUsingLDAP(loginReq.Username, loginReq.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Authentication failed",
			"data":    err.Error(),
		})
	}

	if !success {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid username or password",
			"data":    nil,
		})
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