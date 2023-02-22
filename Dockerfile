FROM golang:alpine as builder

COPY . /app
WORKDIR /app

RUN apk add build-base

#ARG PB_ADMIN_USER
#ENV PB_ADMIN_USER=$PB_ADMIN_USER
#
#ARG PB_ADMIN_PASSWORD
#ENV PB_ADMIN_PASSWORD=$PB_ADMIN_PASSWORD

RUN go mod download

RUN CGO_ENABLED=0 GOOS="linux" GOARCH="amd64" go build -o pocketbase ./main.go

## Deploy
FROM alpine:latest

ARG PB_ADMIN_USER
ENV PB_ADMIN_USER=$PB_ADMIN_USER

ARG PB_ADMIN_PASSWORD
ENV PB_ADMIN_PASSWORD=$PB_ADMIN_PASSWORD

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY --from=builder /app/pocketbase /usr/local/bin/pocketbase

EXPOSE 8090
ENTRYPOINT ["/usr/local/bin/pocketbase", "serve", "--http=0.0.0.0:8090", "--dir=/pb_data"]