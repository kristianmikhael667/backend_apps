package routes

import (
	"backend_apps/controllers"
	"backend_apps/database"
	util "backend_apps/package"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	jwt := jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		ErrorHandler: unauthorize,
	})

	app.Get("/status", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(map[string]string{"status": "OKE"})
	})

	baseUrl := util.Getenv("MIDDLE_URL", "")
	// Without JWT
	app.Post(baseUrl+"/admin-login", controllers.AdminLogin) // done
	// Inisialisasi WebSocket
	controllers.SetupWebSocket(app)

	// Middleware untuk WebSocket
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", database.GetConnection())
		c.Next()
		return nil
	})

	// With JWT Provider
	app.Use(jwt).Post(baseUrl+"/create-provider", controllers.CreateProvider)
	app.Use(jwt).Get(baseUrl+"/get-provider", controllers.GetProvider)
	app.Use(jwt).Get(baseUrl+"/get-provider/:id", controllers.GetDetailProvider)
	app.Use(jwt).Put(baseUrl+"/update-provider/:id", controllers.UpdateProvider)
	app.Use(jwt).Delete(baseUrl+"/delete-provider/:id", controllers.DeleteProvider)

	// With JWT Contact

	app.Use(jwt).Post(baseUrl+"/create-contact", controllers.CreateContact)
	app.Use(jwt).Get(baseUrl+"/get-contact", controllers.GetContact)
	app.Use(jwt).Get(baseUrl+"/get-contact/:id", controllers.GetDetailContact)
	app.Use(jwt).Put(baseUrl+"/update-contact/:id", controllers.UpdateContact)
	app.Use(jwt).Delete(baseUrl+"/delete-contact/:id", controllers.DeleteContact)
}

func unauthorize(c *fiber.Ctx, err error) error {
	return c.Status(401).JSON(fiber.Map{
		"status":  "error",
		"message": err.Error()})
}
