lint: fmt
	golangci-lint run --enable-all --fix

fmt:
	gofmt -w .

gen:
	go generate ./...

test:
	go clean -testcache
	CONFIGOR_ENV=local ROOT_DIR=${PWD} go test -failfast ./...

build:
	docker build . -t soulgarden/swup:pass-0.0.5 --platform linux/amd64
	docker push soulgarden/swup:pass-0.0.5
