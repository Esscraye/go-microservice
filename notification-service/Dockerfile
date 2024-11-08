FROM golang:1.21-alpine
RUN apk --no-cache add curl
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o notification-service .
EXPOSE 8085
CMD ["./notification-service"]
