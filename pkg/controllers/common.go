package controllers

import (
	"github.com/mixedmachine/simple-signin-backend/pkg/repository"
	"github.com/mixedmachine/simple-signin-backend/pkg/util"

	"log"

	"github.com/gofiber/fiber/v2"
)

func AuthRequest(ctx *fiber.Ctx, tokensRepo repository.TokenRepository) (string, error) {
	token := string(ctx.Request().Header.Peek("Authorization"))
	// log.Printf("Token: %s\n", token)
	if token == "" {
		return "", util.ErrInvalidAuthToken
	}
	user, err := tokensRepo.Retrieve(token)
	if user == "" || err != nil {
		log.Printf("User: %s\n", user)
		log.Printf("Error: %s\n", err)
		return "", util.ErrUnauthorized
	}

	return user, nil
}

// func AuthRefreshRequest(ctx *fiber.Ctx, newToken, user string) {
// 	oldToken := ctx.Params("token")
// 	repository.AddToken(newToken, user)
// 	repository.RemoveToken(oldToken)
// }
