build:
	go build ./...

run:
	go run .

test:
	go test ./... -v

check-install-swagger:
	which ../../bin/swagger || GO111MODULE=on go install github.com/go-swagger/go-swagger/cmd/swagger@latest

swagger: check-install-swagger
	GO111MODULE=off ../../bin/swagger generate spec -o ./swagger.yaml --scan-models

check-lint:
	which ../../bin/golint || GO111MODULE=on go install golang.org/x/lint/golint/golint@latest

lint:
	GO111MODULE=off ../../bin/golint ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

build-image:
	docker build --tag storeapi --no-cache=true .

run-docker:
	docker run -d -p 8080:8080 storeapi	