# Building the binary of the App
FROM golang:1.19 AS build

# `boilerplate` should be replaced with your project name
WORKDIR /go/src/deeincom

# Copy all the Code and stuff to compile everything
COPY . .

<<<<<<< HEAD
RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".web.toml"]
=======
# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN go mod download

# Builds the application as a staticly linked one, to allow it to run on alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .


# Moving the binary to the 'final Image' to make it smaller
FROM alpine:latest as release

WORKDIR /app

# Create the `public` dir and copy all the assets into it
RUN mkdir ./static
COPY ./static ./static

# `boilerplate` should be replaced here as well
COPY --from=build /go/src/deeincom/app .

# Add packages
RUN apk -U upgrade \
    && apk add --no-cache dumb-init ca-certificates \
    && chmod +x /app/app

# Exposes port 3000 because our program listens on that port
EXPOSE 3000

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
>>>>>>> 5ed448758ec77912d8f15b1cd516eb78d5d5ec71
