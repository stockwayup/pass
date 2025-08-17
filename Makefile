lint: goimports fmt
	golangci-lint run --fix

fmt:
	gofmt -w .

goimports:
	goimports -w .

gen:
	go generate ./...

test:
	go clean -testcache
	CONFIGOR_ENV=local ROOT_DIR=${PWD} go test -failfast ./...

IMAGE ?= soulgarden/swup:pass-0.0.13
PLATFORM ?= linux/amd64

docker-build:
	docker build . -t $(IMAGE) --platform $(PLATFORM)

docker-push:
	docker push $(IMAGE)

build: docker-build docker-push
