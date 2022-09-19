package controllers

import (
	"github.com/mixedmachine/SimpleAuthBackend/pkg/models"
	"github.com/mixedmachine/SimpleAuthBackend/pkg/repository"
	"github.com/mixedmachine/SimpleAuthBackend/pkg/security"
	"github.com/mixedmachine/SimpleAuthBackend/pkg/util"

	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/asaskevich/govalidator.v9"
)

// UserController defines the interface for user controller
type UserController interface {
	GetUser(ctx *fiber.Ctx) error
	GetUsers(ctx *fiber.Ctx) error
	PutUser(ctx *fiber.Ctx) error
	DeleteUser(ctx *fiber.Ctx) error
}

// userController implements UserController
type userController struct {
	usersRepo  repository.UsersRepository
	tokensRepo repository.TokenRepository
	logger     *log.Logger
}

// NewUserController constructs a new instance of UserController with given repository dependencies
func NewUserController(repos map[string]interface{}) UserController {
	return &userController{
		usersRepo:  repos["users"].(repository.UsersRepository),
		tokensRepo: repos["tokens"].(repository.TokenRepository),
		logger:     util.InitLogger(),
	}
}

/********************************************************
 *				Handler Functions for Users				*
 ********************************************************/

// GetUser returns a user by id
func (c *userController) GetUser(ctx *fiber.Ctx) error {
	c.logger.Debug("GetUser called")
	userId, err := AuthRequest(ctx, c.tokensRepo)
	if err != nil {
		c.logger.Error(err)
		return ctx.
			Status(http.StatusUnauthorized).
			JSON(util.NewJError(err))
	}
	user, err := c.usersRepo.GetById(userId)
	if err != nil {
		c.logger.Error(err)
		return ctx.
			Status(http.StatusInternalServerError).
			JSON(util.NewJError(err))
	}
	c.logger.Debug("GetUser returning user")
	return ctx.
		Status(http.StatusOK).
		JSON(user)
}

// GetUsers returns all users
func (c *userController) GetUsers(ctx *fiber.Ctx) error {
	c.logger.Debug("GetUsers called")
	users, err := c.usersRepo.GetAll()
	if err != nil {
		c.logger.Error(err)
		return ctx.
			Status(http.StatusInternalServerError).
			JSON(util.NewJError(err))
	}
	c.logger.Debug("GetUsers returning users")
	return ctx.
		Status(http.StatusOK).
		JSON(users)
}

// PutUser updates a user by id
func (c *userController) PutUser(ctx *fiber.Ctx) error {
	c.logger.Debug("PutUser called")
	userId, err := AuthRequest(ctx, c.tokensRepo)
	if err != nil {
		c.logger.Error(err)
		return ctx.
			Status(http.StatusUnauthorized).
			JSON(util.NewJError(err))
	}
	var update models.User
	err = ctx.BodyParser(&update)
	if err != nil {
		c.logger.Error(err)
		return ctx.
			Status(http.StatusUnprocessableEntity).
			JSON(util.NewJError(err))
	}
	update.Email = util.NormalizeEmail(update.Email)
	if !govalidator.IsEmail(update.Email) {
		c.logger.Error("Invalid email")
		return ctx.
			Status(http.StatusBadRequest).
			JSON(util.NewJError(util.ErrInvalidEmail))
	}
	exists, err := c.usersRepo.GetByEmail(update.Email)
	if err == mongo.ErrNoDocuments || exists.Id.Hex() == userId {
		c.logger.Debug("Email is unique")
		user, err := c.usersRepo.GetById(userId)
		if err != nil {
			c.logger.Error(err)
			return ctx.
				Status(http.StatusBadRequest).
				JSON(util.NewJError(err))
		}
		if update.Email != "" {
			c.logger.Debug("Updating email")
			user.Email = update.Email
		}
		if update.Password != "" {
			c.logger.Debug("Updating password")
			update.Password, err = security.EncryptPassword(update.Password)
			if err != nil {
				return ctx.
					Status(http.StatusBadRequest).
					JSON(util.NewJError(err))
			}
			user.Password = update.Password
		}
		user.UpdatedAt = time.Now()
		err = c.usersRepo.Update(user)
		if err != nil {
			c.logger.Error(err)
			return ctx.
				Status(http.StatusUnprocessableEntity).
				JSON(util.NewJError(err))
		}
		c.logger.Debug("PutUser returning user")
		return ctx.
			Status(http.StatusOK).
			JSON(user)
	}

	if exists != nil {
		c.logger.Error("Email already exists")
		err = util.ErrEmailAlreadyExists
	}
	c.logger.Error(err)
	return ctx.
		Status(http.StatusBadRequest).
		JSON(util.NewJError(err))
}

// DeleteUser deletes a user by id
func (c *userController) DeleteUser(ctx *fiber.Ctx) error {
	c.logger.Debug("DeleteUser called")
	userId, err := AuthRequest(ctx, c.tokensRepo)
	if err != nil {
		c.logger.Error(err)
		return ctx.
			Status(http.StatusUnauthorized).
			JSON(util.NewJError(err))
	}
	err = c.usersRepo.Delete(userId)
	if err != nil {
		c.logger.Error(err)
		return ctx.
			Status(http.StatusInternalServerError).
			JSON(util.NewJError(err))
	}
	err = c.tokensRepo.Delete(string(ctx.Request().Header.Peek("Authorization")))
	if err != nil {
		c.logger.Error(err)
		return ctx.
			Status(http.StatusInternalServerError).
			JSON(util.NewJError(err))
	}
	c.logger.Debug("DeleteUser returning user")
	ctx.Set("Entity", userId)
	return ctx.SendStatus(http.StatusNoContent)
}
