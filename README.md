<<<<<<< HEAD
# deeincom

Mọi người sửa UI nhớ build css: `cd node && npm run build`. Nếu ko sẽ thiếu style :)).


# requirement

golang 1.16 hoặc mới hơn
postgresql 12 hoặc mới hơn

# chạy

```
go run cmd/deein/main.com web
```

# Sử dụng docker để run app thì:
+ Install docker

```
cd ./project
docker compose up
```
=======
# Deeincom

## Development

### Start the application 

```bash
go run main.go serve
```

### Use local container

```
# Clean packages
make clean-packages

# Generate go.mod & go.sum files
make requirements

# Generate docker image
make build

# Generate docker image with no cache
make build-no-cache

# Run the projec in a local container
make up

# Run local container in background
make up-silent

# Run local container in background with prefork
make up-silent-prefork

# Stop container
make stop

# Start container
make start
```

## Environments
```
APP_ADDR=0.0.0.0:3000
APP_ENV=local
DB_URI=postgresql://user:pwd@localhost:5432/dbName?sslmode=disable
AUTH_SECRET=
AUTH_EXPIRE_DURATION=1h
SMS_KEY=
SMS_SECRET_KEY=
SMS_BRAND_NAME=
```

## Production

```bash
docker build -t gofiber .
docker run -d -p 3000:3000 gofiber ./app -prod
```

Go to http://localhost:3000:
>>>>>>> 5ed448758ec77912d8f15b1cd516eb78d5d5ec71
