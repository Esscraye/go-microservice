FROM golang:1.21-alpine
RUN apk --no-cache add curl
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o user-service .
EXPOSE 8082
CMD ["./user-service"]

