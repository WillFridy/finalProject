FROM golang:1.17-alpine AS build
EXPOSE 8000

WORKDIR /finalProject/
COPY . .

RUN go mod download


CMD [ "go", "run", "main.go" ]