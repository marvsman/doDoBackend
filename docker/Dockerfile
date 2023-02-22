FROM golang:alpine as builder

COPY . /app
WORKDIR /app

RUN apk add build-base

RUN go mod download
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS="linux" GOARCH="amd64" go build -o pocketbase ./main.go

## Deploy
FROM alpine:latest

# get git actions secrets
ARG PB_ADMIN_USER
ENV PB_ADMIN_USER=$PB_ADMIN_USER

ARG PB_ADMIN_PASSWORD
ENV PB_ADMIN_PASSWORD=$PB_ADMIN_PASSWORD

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY --from=builder /app/pocketbase /usr/local/bin/pocketbase

EXPOSE 8090
ENTRYPOINT ["/usr/local/bin/pocketbase", "serve", "--http=0.0.0.0:8090", "--dir=/pb_data"]