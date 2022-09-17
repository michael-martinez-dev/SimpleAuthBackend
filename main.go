package main

import "github.com/mixedmachine/simple-signin-backend/cmd/v1/api"

func main() {
	api.Init()
	api.RunUserAuthApiServer()
}
