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

docker.build:
	docker build -t $(APP_NAME):$(APP_VERSION) .

docker.network.add:
	@if docker network inspect net-$(APP_NAME) >/dev/null 2>&1; then \
		echo "Network net-$(APP_NAME) exists, removing..."; \
		docker network rm net-$(APP_NAME); \
	else \
		echo "Network net-$(APP_NAME) does not exist, creating..."; \
	fi
	docker network create net-$(APP_NAME)

docker.run: swag docker.network.add docker.build
	docker run -d --name $(APP_NAME) --network net-$(APP_NAME) -p 3000:3000 $(APP_NAME):$(APP_VERSION)
