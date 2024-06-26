# Build stage
FROM golang:latest AS build
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o main .

# Run stage
FROM debian:latest
WORKDIR /app
COPY --from=build /app /app

# Uygulamanın çalışması için gerekli bağımlılıkları ekleyin
RUN apt-get update && apt-get install -y ca-certificates

# ENTRYPOINT ve EXPOSE direktifleri
EXPOSE 8080
ENTRYPOINT ["./main"]
