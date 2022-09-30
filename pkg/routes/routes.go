package routes

import (
	"github.com/mixedmachine/SimpleAuthBackend/pkg/controllers"

	"fmt"

	"github.com/gofiber/fiber/v2"
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
	app.Get("/ping", r.authController.Ping)

	/********************************
	 * Authentication & User Routes *
	 ********************************/
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
