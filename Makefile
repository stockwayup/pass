lint: fmt
	golangci-lint run --enable-all --fix

fmt:
	gofmt -w .

gen:
	go generate ./...

build:
	docker build . -t soulgarden/swup:pass-0.0.1 --platform linux/amd64
	docker push soulgarden/swup:pass-0.0.1
