#!/bin/bash

# Fonction pour obtenir un token JWT
get_jwt_token() {
    local user_id=$1
    local password=$2
    local response=$(curl -s -X POST -H "Content-Type: application/json" -d "{\"user_id\":\"$user_id\",\"password\":\"$password\"}" http://localhost:8080/login)
    echo $(echo $response | jq -r .token)
}

# Fonction pour afficher un séparateur
print_separator() {
    echo "========================================"
}

user_service_tests() {
    local port=8082
    local uri=http://localhost:$port

    print_separator
    echo "Testing user-service on port $port"
    print_separator

    echo "1. Reading all users"
    curl -s -H "Authorization: $jwt_token" $uri/users
    echo

    echo "2. Creating a user"
    curl -s -X POST -H "Content-Type: application/json" -H "Authorization: $jwt_token" -d '{"id":"3","name":"User3","email":"user3@example.com","password":"password"}' $uri/users
    echo

    echo "3. Reading all users"
    curl -s -H "Authorization: $jwt_token" $uri/users
    echo

    echo "4. Updating a user"
    curl -s -X PUT -H "Content-Type: application/json" -H "Authorization: $jwt_token" -d '{"id":"3","name":"UpdatedUser","email":"updateduser@example.com","password":"newpassword"}' $uri/users/3
    echo

    echo "5. Reading the updated user"
    curl -s -H "Authorization: $jwt_token" $uri/users/3
    echo

    echo "6. Deleting a user"
    curl -s -X DELETE -H "Authorization: $jwt_token" $uri/users/3
    echo

    echo "7. Reading all users"
    curl -s -H "Authorization: $jwt_token" $uri/users
    echo
}

product_service_tests() {
    local port=8081
    local uri=http://localhost:$port

    print_separator
    echo "Testing product-service on port $port"
    print_separator

    echo "1. Reading all products"
    curl -s -H "Authorization: $jwt_token" $uri/products
    echo

    echo "2. Creating a product"
    curl -s -X POST -H "Content-Type: application/json" -H "Authorization: $jwt_token" -d '{"id":"3","name":"Product3","category":"Category3","price":300.0}' $uri/products
    echo

    echo "3. Reading all products"
    curl -s -H "Authorization: $jwt_token" $uri/products
    echo

    echo "4. Updating a product"
    curl -s -X PUT -H "Content-Type: application/json" -H "Authorization: $jwt_token" -d '{"id":"3","name":"UpdatedProduct","category":"UpdatedCategory","price":350.0}' $uri/products/3
    echo

    echo "5. Reading the updated product"
    curl -s -H "Authorization: $jwt_token" $uri/products/3
    echo

    echo "6. Deleting a product"
    curl -s -X DELETE -H "Authorization: $jwt_token" $uri/products/3
    echo

    echo "7. Reading all products"
    curl -s -H "Authorization: $jwt_token" $uri/products
    echo
}

order_service_tests() {
    local port=8083
    local uri=http://localhost:$port

    print_separator
    echo "Testing order-service on port $port"
    print_separator

    echo "1. Reading all orders"
    curl -s -H "Authorization: $jwt_token" $uri/orders
    echo

    echo "2. Creating many orders"
    curl -s -X POST -H "Content-Type: application/json" -H "Authorization: $jwt_token" -d '{"id":"1","user_id":"1","product_id":"1","quantity":2,"status":"pending"}' $uri/orders
    curl -s -X POST -H "Content-Type: application/json" -H "Authorization: $jwt_token" -d '{"id":"2","user_id":"1","product_id":"2","quantity":2,"status":"pending"}' $uri/orders
    curl -s -X POST -H "Content-Type: application/json" -H "Authorization: $jwt_token" -d '{"id":"3","user_id":"2","product_id":"1","quantity":2,"status":"pending"}' $uri/orders
    echo

    echo "3. Reading all orders"
    curl -s -H "Authorization: $jwt_token" $uri/orders
    echo

    echo "4. Updating an order"
    curl -s -X PUT -H "Content-Type: application/json" -H "Authorization: $jwt_token" -d '{"id":"1","user_id":"1","product_id":"1","quantity":3,"status":"shipped"}' $uri/orders/1
    echo

    echo "5. Reading the updated order"
    curl -s -H "Authorization: $jwt_token" $uri/orders/1
    echo

    echo "6. Deleting an order"
    curl -s -X DELETE -H "Authorization: $jwt_token" $uri/orders/1
    echo

    echo "7. Reading all orders"
    curl -s -H "Authorization: $jwt_token" $uri/orders
    echo
}

payment_service_tests() {
    local port=8084
    local uri=http://localhost:$port

    print_separator
    echo "Testing payment-service on port $port"
    print_separator

    echo "1. Reading all payments"
    curl -s -H "Authorization: $jwt_token" $uri/payments
    echo

    echo "2. Creating a payment"
    curl -s -X POST -H "Content-Type: application/json" -H "Authorization: $jwt_token" -d '{"id":"1","order_id":"2","amount":100.0,"status":"pending"}' $uri/payments
    echo

    echo "3. Reading all payments"
    curl -s -H "Authorization: $jwt_token" $uri/payments
    echo

    echo "4. Updating a payment"
    curl -s -X PUT -H "Content-Type: application/json" -H "Authorization: $jwt_token" -d '{"id":"1","order_id":"2","amount":150.0,"status":"completed"}' $uri/payments/1
    echo

    echo "5. Reading the updated payment"
    curl -s -H "Authorization: $jwt_token" $uri/payments/1
    echo

    echo "6. Deleting a payment"
    curl -s -X DELETE -H "Authorization: $jwt_token" $uri/payments/1
    echo

    echo "7. Reading all payments"
    curl -s -H "Authorization: $jwt_token" $uri/payments
    echo
}

notification_service_tests() {
    local port=8085
    local uri=http://localhost:$port

    print_separator
    echo "Testing notification-service on port $port"
    print_separator

    echo "1. Reading all notifications"
    curl -s -H "Authorization: $jwt_token" $uri/notifications
    echo

    echo "2. Creating a notification"
    curl -s -X POST -H "Content-Type: application/json" -H "Authorization: $jwt_token" -d '{"id":"1","user_id":"1","message":"Your order has been shipped","status":"pending"}' $uri/notifications
    echo

    echo "3. Reading all notifications"
    curl -s -H "Authorization: $jwt_token" $uri/notifications
    echo

    echo "4. Updating a notification"
    curl -s -X PUT -H "Content-Type: application/json" -H "Authorization: $jwt_token" -d '{"id":"1","user_id":"1","message":"Your order has been delivered","status":"completed"}' $uri/notifications/1
    echo

    echo "5. Reading the updated notification"
    curl -s -H "Authorization: $jwt_token" $uri/notifications/1
    echo

    echo "6. Deleting a notification"
    curl -s -X DELETE -H "Authorization: $jwt_token" $uri/notifications/1
    echo

    echo "7. Reading all notifications"
    curl -s -H "Authorization: $jwt_token" $uri/notifications
    echo
}

# Obtenir un token JWT pour les tests
jwt_token=$(get_jwt_token "user1@example.com" "password")
echo "JWT token: $jwt_token"
print_separator

# Tester les services
user_service_tests
product_service_tests
order_service_tests
payment_service_tests
notification_service_tests