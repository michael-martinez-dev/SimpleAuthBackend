package api

import (
	"fmt"
	"os"

	"github.com/mixedmachine/SimpleAuthBackend/pkg/controllers"
	"github.com/mixedmachine/SimpleAuthBackend/pkg/db"
	"github.com/mixedmachine/SimpleAuthBackend/pkg/repository"
	"github.com/mixedmachine/SimpleAuthBackend/pkg/routes"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Panicln(err)
	}
}

func RunUserAuthApiServer() {
	mConn := db.NewMongoConnection()
	rConn := db.NewRedisConnection()
	defer mConn.Close()
	defer rConn.Close()

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	userRepo := repository.NewUserRepository(mConn)
	tokenRepo := repository.NewTokenRepository(rConn)
	repos := map[string]interface{}{
		"users":  userRepo,
		"tokens": tokenRepo,
	}
	authController := controllers.NewAuthController(repos)
	userController := controllers.NewUserController(repos)

	authRoutes := routes.NewAuthRoutes(authController, userController)
	authRoutes.Install(app)

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}
	addr := os.Getenv("API_ADDR")
	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", addr, port)))
}
