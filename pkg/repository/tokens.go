package repository

import (
	"github.com/mixedmachine/simple-signin-backend/pkg/db"
	"github.com/mixedmachine/simple-signin-backend/pkg/util"

	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-redis/redis"
)

const experationTime = 15 // minutes

// TokenRepository is an interface for token repository
type TokenRepository interface {
	Create(token, user string, expire bool) error
	Retrieve(token string) (string, error)
	Delete(token string) error
}

// tokensRepository is a struct for token repository
type tokensRepository struct {
	rClient *redis.Client
	rCtx    context.Context
	logger  *log.Logger
}

// NewTokenRepository returns a new token repository with redis connection
func NewTokenRepository(conn db.RedisConnection) TokenRepository {
	return &tokensRepository{
		rClient: conn.GetClient(),
		rCtx:    conn.GetContext(),
		logger:  util.InitLogger(),
	}
}

// Create creates a new token for user in the database with an option to expire
// the token after a certain amount of time
func (r *tokensRepository) Create(token, user string, expire bool) error {
	exp := 0
	exp_str := "never"

	if expire {
		exp = experationTime
		exp_str = fmt.Sprintf("in %d minutes", experationTime)
	}

	err := r.rClient.Set(
		token, user,
		time.Duration(exp)*time.Minute,
	).Err()

	if err != nil {
		return err
	}
	r.logger.Debugf("Created token for user %s that expires %s\n", user, exp_str)

	return nil
}

// Retrieve retrieves a user from the database by token
func (r *tokensRepository) Retrieve(token string) (string, error) {
	userId, err := r.rClient.Get(token).Result()
	if err != nil {
		return "", err
	}
	r.logger.Debugf("Retrieved token for user %s\n", userId)
	return userId, nil
}

// Delete deletes a token from the database
func (r *tokensRepository) Delete(token string) error {
	err := r.rClient.Del(token).Err()
	if err != nil {
		return err
	}
	r.logger.Debugf("Deleted token\n")
	return nil
}
