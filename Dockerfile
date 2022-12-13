# Building the binary of the App
FROM golang:1.19 AS build

# `boilerplate` should be replaced with your project name
WORKDIR /go/src/deeincom

# Copy all the Code and stuff to compile everything
COPY . .

RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".web.toml"]