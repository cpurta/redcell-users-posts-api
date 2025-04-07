FROM golang:1.22-alpine3.19 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./publish/bin/users-posts-api /app/cmd/server/main.go

FROM alpine:3.19

RUN apk add curl

COPY --from=build /app/publish/bin .

ENTRYPOINT [ "./users-posts-api", "start" ]