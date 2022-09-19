package test_util

import (
	"github.com/mixedmachine/SimpleAuthBackend/pkg/util"

	"testing"
)

func TestInitLogger(t *testing.T) {
	logger := util.InitLogger()
	if logger == nil {
		t.Error("logger is nil")
	}
}

func TestNewJError(t *testing.T) {
	jerr := util.NewJError(nil)
	if jerr.Error != "generic error" {
		t.Error("jerr.Error is not generic error")
	}
}

func TestNormalizeEmail(t *testing.T) {
	email := util.NormalizeEmail(" tester@email.com ")
	if email != "tester@email.com" {
		t.Error("email is not normalized")
	}
}
