FROM golang:1.21-alpine
RUN apk --no-cache add curl
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o payment-service .
EXPOSE 8084
CMD ["./payment-service"]

