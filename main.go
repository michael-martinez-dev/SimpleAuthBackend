package main

import "github.com/mixedmachine/SimpleAuthBackend/cmd/v1/api"

func main() {
	api.Init()
	api.RunUserAuthApiServer()
}
