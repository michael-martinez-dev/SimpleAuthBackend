package controllers

import (
	"github.com/mixedmachine/SimpleAuthBackend/pkg/models"
	"github.com/mixedmachine/SimpleAuthBackend/pkg/repository"
	"github.com/mixedmachine/SimpleAuthBackend/pkg/security"
	"github.com/mixedmachine/SimpleAuthBackend/pkg/util"

	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/asaskevich/govalidator.v9"
)

// AuthController interface defines the contract for the AuthController
type AuthController interface {
	Ping(ctx *fiber.Ctx) error
	SignUp(ctx *fiber.Ctx) error
	SignIn(ctx *fiber.Ctx) error
	RefreshToken(ctx *fiber.Ctx) error
	Authenticator(ctx *fiber.Ctx) error
}

// authController struct implements the AuthController interface
type authController struct {
	usersRepo  repository.UsersRepository
	tokensRepo repository.TokenRepository
	logger     *log.Logger
}

// NewAuthController constructs a new instance of AuthController with given repository dependencies
func NewAuthController(repos map[string]interface{}) AuthController {
	return &authController{
		usersRepo:  repos["users"].(repository.UsersRepository),
		tokensRepo: repos["tokens"].(repository.TokenRepository),
		logger:     util.InitLogger(),
	}
}

// Ping Handler Function for Health Check
func (c *authController) Ping(ctx *fiber.Ctx) error {
	return ctx.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"message": "pong",
		})
}

/********************************************************
 * 		  Handler Functions for Authentication			*
 ********************************************************/

// SignUp Handler Function verifies the user input and creates a new user in the database
func (c *authController) SignUp(ctx *fiber.Ctx) error {
	c.logger.Debug("SignUp Handler Function")
	var newUser models.User

	err := ctx.BodyParser(&newUser)
	if err != nil {
		c.logger.Errorf("SignUp Handler Function| Error parsing request body: %v", err)
		return ctx.
			Status(http.StatusUnprocessableEntity).
			JSON(util.NewJError(err))
	}

	err = verifyUser(&newUser, c)
	if err != nil {
		c.logger.Errorf("SignUp Handler Function| Error verifying user: %v", err)
		return ctx.
			Status(http.StatusBadRequest).
			JSON(util.NewJError(err))
	}

	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = newUser.CreatedAt
	newUser.Id = primitive.NewObjectID()

	err = c.usersRepo.Save(&newUser)
	if err != nil {
		c.logger.Errorf("SignUp Handler Function| Error saving user: %v", err)
		return ctx.
			Status(http.StatusBadRequest).
			JSON(util.NewJError(err))
	}
	c.logger.Debugf("SignUp Handler Function| User saved: %v", newUser)
	return ctx.
		Status(http.StatusCreated).
		JSON(newUser)
}

// SignIn Handler Function verifies the user input and returns a new token
func (c *authController) SignIn(ctx *fiber.Ctx) error {
	c.logger.Debug("SignIn Handler Function")
	var input models.User
	err := ctx.BodyParser(&input)
	if err != nil {
		c.logger.Errorf("SignIn Handler Function| Error parsing request body: %v", err)
		return ctx.
			Status(http.StatusUnprocessableEntity).
			JSON(util.NewJError(err))
	}

	input.Email = util.NormalizeEmail(input.Email)
	user, err := c.usersRepo.GetByEmail(input.Email)
	if err != nil {
		c.logger.Errorf("c.usersRepo.GetByEmail| %s signin failed: %v\n", input.Email, err.Error())
		return ctx.
			Status(http.StatusUnauthorized).
			JSON(util.NewJError(util.ErrInvalidCredentials))
	}

	err = security.VerifyPassword(user.Password, input.Password)
	if err != nil {
		c.logger.Errorf("security.VerifyPassword| %s signin failed: %v\n", input.Email, err.Error())
		return ctx.
			Status(http.StatusUnauthorized).
			JSON(util.NewJError(util.ErrInvalidCredentials))
	}

	token, err := security.NewToken(user.Id.Hex())
	if err != nil {
		c.logger.Errorf("security.NewToken| %s signin failed: %v\n", input.Email, err.Error())
		return ctx.
			Status(http.StatusUnauthorized).
			JSON(util.NewJError(err))
	}
	c.tokensRepo.Create(token, user.Id.Hex(), true)
	c.logger.Debugf("SignIn Handler Function| Token created: %v", token)
	return ctx.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"user":  user,
			"token": token,
		})
}

// RefreshToken Handler Function verifies the user input removes old token and returns a new token
func (c *authController) RefreshToken(ctx *fiber.Ctx) error {
	c.logger.Debug("RefreshToken Handler Function")
	userId, err := AuthRequest(ctx, c.tokensRepo)
	if err != nil {
		c.logger.Errorf("AuthRequest| %s refresh failed: %v\n", userId, err.Error())
		return ctx.
			Status(http.StatusUnauthorized).
			JSON(util.NewJError(err))
	}

	token, err := security.NewToken(userId)
	if err != nil {
		c.logger.Errorf("security.NewToken| %s refresh failed: %v\n", userId, err.Error())
		return ctx.
			Status(http.StatusUnauthorized).
			JSON(util.NewJError(err))
	}

	c.tokensRepo.Create(token, userId, true)
	c.tokensRepo.Delete(string(ctx.Request().Header.Peek("Authorization")))

	c.logger.Debugf("RefreshToken Handler Function| Token created: %v", token)
	return ctx.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"token": token,
		})
}

// Authenticator Handler Function takes the token from the request header and returns the user id
// associated with the token
func (c *authController) Authenticator(ctx *fiber.Ctx) error {
	c.logger.Debug("Authenticator Handler Function")
	userId, err := AuthRequest(ctx, c.tokensRepo)
	if err != nil {
		c.logger.Errorf("AuthRequest| %s authentication failed: %v\n", userId, err.Error())
		return ctx.
			Status(http.StatusUnauthorized).
			JSON(util.NewJError(err))
	}
	c.logger.Debugf("Authenticator Handler Function| User id: %v", userId)
	return ctx.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"user_id": userId,
		})
}

/********************************************************
* 					Helper functions					*
*********************************************************/

// verifyUser verifies the user input and returns an error if the input is invalid
func verifyUser(user *models.User, c *authController) error {
	c.logger.Debug("verifyUser Function")
	if user.Email == "" {
		return util.ErrInvalidEmail
	}
	if user.Password == "" {
		return util.ErrEmptyPassword
	}

	user.Email = util.NormalizeEmail(user.Email)
	if !govalidator.IsEmail(user.Email) {
		return util.ErrInvalidEmail
	}

	exists, err := c.usersRepo.GetByEmail(user.Email)
	if err != mongo.ErrNoDocuments {
		return err
	}

	if exists != nil {
		err = util.ErrEmailAlreadyExists
		return err
	}

	if strings.TrimSpace(user.Password) == "" {
		return util.ErrEmptyPassword
	}

	user.Password, err = security.EncryptPassword(user.Password)
	if err != nil {
		return err
	}
	c.logger.Debug("verifyUser Function| User verified")
	return nil
}
