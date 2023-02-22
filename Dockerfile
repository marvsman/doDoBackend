FROM golang:1.18.3-alpine as builder

WORKDIR /app

RUN export VERSION=$(git-semver-describe --tags)

ARG PB_ADMIN_USER
ENV PB_ADMIN_USER=$PB_ADMIN_USER

ARG PB_ADMIN_PASSWORD
ENV PB_ADMIN_PASSWORD=$PB_ADMIN_PASSWORD

COPY * ./
RUN go mod download

RUN GOOS="linux" GOARCH="amd64" go build -o /pocketbase

## Deploy
FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY --from=builder /app/pocketbase /usr/local/bin/pocketbase

EXPOSE 8090
ENTRYPOINT ["/usr/local/bin/pocketbase", "serve", "--http=0.0.0.0:8090", "--dir=/pb_data", "--publicDir=/pb_public"]