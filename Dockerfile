FROM golang:1.23.3-alpine3.20 as builder

ENV GOPATH=/go

COPY . $GOPATH/src/github.com/stockwayup/pass

WORKDIR $GOPATH/src/github.com/stockwayup/pass

RUN go get -u -t github.com/tinylib/msgp && \
    go install github.com/tinylib/msgp && \
    go generate ./... && \
    go build -o /bin/stockwayup

FROM alpine:3.20

RUN adduser -S www-data -G www-data

COPY --from=builder --chown=www-data /bin/stockwayup /bin/stockwayup

RUN chmod +x /bin/stockwayup

USER www-data

CMD ["/bin/stockwayup"]