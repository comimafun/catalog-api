## Build
FROM golang:1.21.10-alpine AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main ./cmd/api/main.go

## Deploy
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/main .
CMD ["./main"]