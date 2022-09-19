package test_models

import (
	"github.com/mixedmachine/SimpleAuthBackend/pkg/models"

	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	t.Run("should create a new user", func(t *testing.T) {
		email := ""
		password := "password"
		user := models.NewUser(email, password)
		if user.Email != email {
			t.Errorf("expected %v, got %v", email, user.Email)
		}
		if user.Password != password {
			t.Errorf("expected %v, got %v", password, user.Password)
		}
		if user.CreatedAt.IsZero() {
			t.Errorf("expected %v, got %v", time.Now(), user.CreatedAt)
		}
		if user.UpdatedAt.IsZero() {
			t.Errorf("expected %v, got %v", time.Now(), user.UpdatedAt)
		}
	})
}
