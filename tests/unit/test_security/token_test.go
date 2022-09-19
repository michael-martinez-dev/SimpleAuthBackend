package test_security

import (
	"github.com/mixedmachine/SimpleAuthBackend/pkg/security"

	"testing"
)

func TestNewToken(t *testing.T) {
	token, err := security.NewToken("test")
	if err != nil {
		t.Error(err)
	}
	if token == "" {
		t.Error("token is empty")
	}
}

func TestParseToken(t *testing.T) {
	token, err := security.NewToken("test")
	if err != nil {
		t.Error(err)
	}
	if token == "" {
		t.Error("token is empty")
	}
	claims, err := security.ParseToken(token)
	if err != nil {
		t.Error(err)
	}
	if claims.Id != "test" {
		t.Error("claims.Id is not test")
	}
}

func TestParseTokenFail(t *testing.T) {
	_, err := security.ParseToken("test")
	if err == nil {
		t.Error("token is invalid")
	}
}
