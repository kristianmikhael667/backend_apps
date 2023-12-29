package controllers

import (
	"backend_apps/database"
	"backend_apps/formatters"
	"backend_apps/models"
	util "backend_apps/package"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

func AdminLogin(c *fiber.Ctx) error {
	payload := new(formatters.RequestBodyAuthAdmin)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	// Check Username
	var users models.User
	if err := database.GetConnection().Model(&models.User{}).Where("email = ? ", payload.Email).First(&users).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Your credential information is invalid."})
		}
		return err
	}

	// Check Password
	isVerify := users.VerifyHash(payload.Password, users.PasswordHash)
	if isVerify == false {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Your credential information is invalid."})
	}

	// Create token jwt
	now := time.Now()
	future := now.Add(2 * 24 * time.Hour)

	jti, _ := util.GenerateBase62EncodedRandomBytes(16)

	secret := os.Getenv("JWT_SECRET")
	signMethodToken := jwt.New(jwt.SigningMethodHS256)
	claims := signMethodToken.Claims.(jwt.MapClaims)
	claims["iat"] = now.Unix()
	claims["exp"] = future.Unix()
	claims["jti"] = jti
	claims["sub"] = users.FullName
	claims["isAdmin"] = true
	claims["uid"] = users.Uid.String()
	token, _ := signMethodToken.SignedString([]byte(secret))

	// Response
	responseAdmin := formatters.ResponseAdminUser{
		Status: "Ok",
		Token:  token,
		User: formatters.User{
			Uid:       users.Uid.String(),
			FullName:  users.FullName,
			Email:     users.Email,
			Status:    users.Status,
			CreatedAt: users.CreatedAt.String(),
			UpdatedAt: users.UpdatedAt.String(),
		},
	}
	return c.Status(http.StatusCreated).JSON(responseAdmin)
}
