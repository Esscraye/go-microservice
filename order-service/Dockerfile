FROM golang:1.21-alpine
RUN apk --no-cache add curl
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o order-service .
EXPOSE 8083
CMD ["./order-service"]

