FROM golang:1.22.2-alpine3.18 as builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go clean --modcache
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o cmd/banners/bin/main ./cmd/banners/

FROM alpine:latest
WORKDIR /banners
COPY /.env .
COPY /config/config.yaml ./config/
COPY --from=builder /app/cmd/banners/bin/main .
CMD ["./main"]