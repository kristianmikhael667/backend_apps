package main

import (
	"backend_apps/database"
	"backend_apps/database/migration"
	"backend_apps/database/seeder"
	"backend_apps/routes"
	"flag"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	database.GetConnection()
}

func main() {
	lumberJack := lumberjack.Logger{
		Filename:   "logs/fiber.log",
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	multi := zerolog.MultiLevelWriter(os.Stdout, &lumberJack)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	database.CreateConnection()

	var migrate string
	var seed string

	flag.StringVar(
		&migrate,
		"migrate",
		"none",
		`this argument for check if user want to migrate table, rollback table, or status migration
to use this flag:
	use -migrate=migrate for migrate table
	use -migrate=rollback for rollback table
	use -migrate=status for get status migration`,
	)

	flag.StringVar(
		&seed,
		"seed",
		"none",
		`this argument for check if user want to seed table
to use this flag:
	use -seed=all to seed all table`,
	)

	flag.Parse()

	if migrate == "migrate" {
		migration.Migrate()
	} else if migrate == "rollback" {
		migration.Rollback()
	} else if migrate == "status" {
		migration.Status()
	} else {
		log.Print("No Key Migrate")
	}

	if seed == "all" {
		seeder.NewSeeder().DeleteAll()
		seeder.NewSeeder().SeedAll()
	}

	app := fiber.New(fiber.Config{
		Prefork: true,
	})

	// logger setting
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
		Output: log.Logger,
	}))

	// setupRoutes(app)
	routes.SetupRoutes(app)
	log.Fatal().AnErr("Fatal: ", app.Listen(":"+os.Getenv("SERVER_PORT")))
}
