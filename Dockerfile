## Build
FROM golang:1.21.10-alpine as build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main ./cmd/api/main.go

## Deploy
FROM golang:1.21.10-alpine
WORKDIR /app
COPY --from=build /app/main .
CMD ["./main"]