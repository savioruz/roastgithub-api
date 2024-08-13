APP_NAME = roastgithub-api
APP_VERSION = 0.0.1
APP_DESCRIPTION = Roasting GitHub API
APP_AUTHOR = savioruz

clean:
	rm -rf ./build

critic:
	gocritic check -enableAll ./...

security:
	gosec ./...

lint:
	golangci-lint run ./...

swag:
	swag init

test: clean critic security lint
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out

build: test
	mkdir -p ./build
	CGO_ENABLED=0 go build -ldflags="-w -s" -o ./build/$(APP_NAME) main.go

run: swag build
	./build/$(APP_NAME)
