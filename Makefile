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

build:
	docker build . -t soulgarden/swup:pass-0.0.11 --platform linux/amd64
	docker push soulgarden/swup:pass-0.0.11
