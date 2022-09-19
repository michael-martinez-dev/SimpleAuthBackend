package routes

import (
	"github.com/mixedmachine/SimpleAuthBackend/pkg/controllers"

	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

const apiVersion = "v1"

type authRoutes struct {
	authController controllers.AuthController
	userController controllers.UserController
}

func NewAuthRoutes(authController controllers.AuthController, userController controllers.UserController) Routes {
	return &authRoutes{
		authController: authController,
		userController: userController,
	}
}

func (r *authRoutes) Install(app *fiber.App) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.
			Status(http.StatusOK).
			JSON(bson.D{
				{Key: "message", Value: "Welcome to EfficientLife"},
				{Key: "service", Value: "user-auth"},
				{Key: "author", Value: "MixedMachine"},
				{Key: "status", Value: http.StatusOK},
				{Key: "version", Value: apiVersion},
				{Key: "api_base_endpoint", Value: "/api/" + apiVersion},
				{Key: "api_endpoints", Value: []string{
					"/ping",
					"/signup",
					"/signin",
					"/refresh",
					"/users/",
					"/users/:id",
					"/auth/:id",
				}},
			})
	})
	api := app.Group(fmt.Sprintf("/api/%s", apiVersion))

	// Health check
	api.Get("/ping", r.authController.Ping)

	// Authentication
	api.Post("/signup", r.authController.SignUp)
	api.Post("/signin", r.authController.SignIn)
	api.Post("/refresh", r.authController.RefreshToken)
	api.Get("/auth", r.authController.Authenticator)

	// Users management
	usersGroup := api.Group("/users")
	usersGroup.Get("/", r.userController.GetUsers)
	usersGroup.Get("/:id", r.userController.GetUser)
	usersGroup.Put("/:id", r.userController.PutUser)
	usersGroup.Delete("/:id", r.userController.DeleteUser)
}
