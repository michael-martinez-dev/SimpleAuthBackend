package test_security

import (
	"github.com/mixedmachine/SimpleAuthBackend/pkg/security"

	"testing"
)

func TestEncryptPassword(t *testing.T) {
	encrypted, err := security.EncryptPassword("password")
	if err != nil {
		t.Error("error encrypting password")
	}
	if len(encrypted) == 0 {
		t.Error("encrypted password is empty")
	}
}

func TestVerifyPassword(t *testing.T) {
	encrypted, err := security.EncryptPassword("password")
	if err != nil {
		t.Error("error encrypting password")
	}
	err = security.VerifyPassword(encrypted, "password")
	if err != nil {
		t.Error("error verifying password")
	}
}

func TestVerifyPasswordFail(t *testing.T) {
	encrypted, err := security.EncryptPassword("password")
	if err != nil {
		t.Error("error encrypting password")
	}
	err = security.VerifyPassword(encrypted, "wrongpassword")
	if err == nil {
		t.Error("error verifying password")
	}
}
