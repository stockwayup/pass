FROM golang:1.25.0-alpine3.21 AS builder

ENV GOPATH=/go

COPY . $GOPATH/src/github.com/stockwayup/pass

WORKDIR $GOPATH/src/github.com/stockwayup/pass

RUN GOFLAGS= go install github.com/tinylib/msgp@v1.3.0 && \
    go generate ./... && \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /bin/stockwayup

FROM alpine:3.21

RUN adduser -S www-data -G www-data

COPY --from=builder --chown=www-data /bin/stockwayup /bin/stockwayup

RUN chmod +x /bin/stockwayup

USER www-data

CMD ["/bin/stockwayup"]
