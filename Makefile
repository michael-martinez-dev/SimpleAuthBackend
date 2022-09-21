APP_NAME = user-auth
APP_VERSION = v1
APP_BIN = server


.PHONY: db dev pipeline image image-push image-run compose 


build: main.go pkg/* cmd/*
	go mod tidy
	go build -o ./bin/$(APP_BIN).exe main.go

db:
	docker compose -f ./build/docker-compose.db.yml up -d

dev: db build 
	./bin/$(APP_BIN).exe

pipeline:
	go fmt ./...
	golangci-lint run

dockerfile:
	go build -o ./bin/$(APP_BIN) main.go

image:
	docker build -f ./build/Dockerfile -t $(APP_NAME):latest .
	docker build -f ./build/Dockerfile -t $(APP_NAME):$(APP_VERSION) .

image-push:
	docker push $(APP_NAME):latest
	docker push $(APP_NAME):$(APP_VERSION)

# This will not be able to connect to mongo unless you change the .env
# to a reachable host. Instead use compose.
image-run: image
	docker run -d \
	-p 8080:8080 \
	--env-file .env \
	--name $(APP_NAME) \
	$(APP_NAME):latest

compose: image
	docker compose -f ./build/docker-compose.db.yml up --build -d
	docker compose -f ./build/docker-compose.api.yml up --build -d

clean:
	rm -f ./bin/$(APP_BIN)
	docker rm -f $(APP_NAME)
	docker compose -f ./build/docker-compose.db.yml down
	docker compose -f ./build/docker-compose.api.yml down
