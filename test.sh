#!/bin/bash

docker compose down --remove-orphans -v --rmi all && docker system prune -af && docker volume prune -af
docker compose up -d --build

# Initialiser la base de données avec des données de test
initialize_db() {
    echo "Initializing database with test data"
    # Ajouter des utilisateurs
    curl -s -X POST -H "Content-Type: application/json" -d '{"id":"1","name":"User1","email":"user1@example.com","password":"password"}' http://localhost:8082/users
    curl -s -X POST -H "Content-Type: application/json" -d '{"id":"2","name":"User2","email":"user2@example.com","password":"password"}' http://localhost:8082/users

    # Ajouter des produits
    curl -s -X POST -H "Content-Type: application/json" -d '{"id":"1","name":"Product1","category":"Category1","price":100.0}' http://localhost:8081/products
    curl -s -X POST -H "Content-Type: application/json" -d '{"id":"2","name":"Product2","category":"Category2","price":200.0}' http://localhost:8081/products
}

# Fonction pour obtenir un token JWT
get_jwt_token() {
    local user_id=$1
    local password=$2
    local response=$(curl -s -X POST -H "Content-Type: application/json" -d "{\"user_id\":\"$user_id\",\"password\":\"$password\"}" http://localhost:8080/login)
    echo $(echo $response | jq -r .token)
}

product_service_tests() {
    local port=8081
    local uri=http://localhost:$port

    echo "Testing product-service on port $port"

    echo "Reading all products"
    curl -s -H "Authorization: Bearer $jwt_token" $uri/products

    echo "Creating a product"
    curl -s -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $jwt_token" -d '{"id":"3","name":"Product3","category":"Category3","price":300.0}' $uri/products

    echo "Reading all products"
    curl -s -H "Authorization: Bearer $jwt_token" $uri/products

    echo "Updating a product"
    curl -s -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer $jwt_token" -d '{"id":"3","name":"UpdatedProduct","category":"UpdatedCategory","price":350.0}' $uri/products/3

    echo "Reading the updated product"
    curl -s -H "Authorization: Bearer $jwt_token" $uri/products/3

    echo "Deleting a product"
    curl -s -X DELETE -H "Authorization: Bearer $jwt_token" $uri/products/3

    echo "Reading all products"
    curl -s -H "Authorization: Bearer $jwt_token" $uri/products
}

user_service_tests() {
    local port=8082
    local uri=http://localhost:$port

    echo "Testing user-service on port $port"

    echo "Reading all users"
    curl -s -H "Authorization: Bearer $jwt_token" $uri/users

    echo "Creating a user"
    curl -s -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $jwt_token" -d '{"id":"3","name":"User3","email":"user3@example.com","password":"password"}' $uri/users

    echo "Reading all users"
    curl -s -H "Authorization: Bearer $jwt_token" $uri/users

    echo "Updating a user"
    curl -s -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer $jwt_token" -d '{"id":"3","name":"UpdatedUser","email":"updateduser@example.com","password":"newpassword"}' $uri/users/3

    echo "Reading the updated user"
    curl -s -H "Authorization: Bearer $jwt_token" $uri/users/3

    echo "Deleting a user"
    curl -s -X DELETE -H "Authorization: Bearer $jwt_token" $uri/users/3

    echo "Reading all users"
    curl -s -H "Authorization: Bearer $jwt_token" $uri/users
}

order_service_tests() {
    local port=8083
    local uri=http://localhost:$port

    echo "Testing order-service on port $port"

    echo "Reading all orders"
    curl -s -H "Authorization: Bearer $jwt_token" $uri/orders

    echo "Creating an order"
    curl -s -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $jwt_token" -d '{"id":"1","user_id":"1","product_id":"1","quantity":2,"status":"pending"}' $uri/orders

    echo "Reading all orders"
    curl -s -H "Authorization: Bearer $jwt_token" $uri/orders

    echo "Updating an order"
    curl -s -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer $jwt_token" -d '{"id":"1","user_id":"1","product_id":"1","quantity":3,"status":"shipped"}' $uri/orders/1

    echo "Reading the updated order"
    curl -s -H "Authorization: Bearer $jwt_token" $uri/orders/1

    echo "Deleting an order"
    curl -s -X DELETE -H "Authorization: Bearer $jwt_token" $uri/orders/1

    echo "Reading all orders"
    curl -s -H "Authorization: Bearer $jwt_token" $uri/orders
}

payment_service_tests() {
    local port=8084
    local uri=http://localhost:$port

    echo "Testing payment-service on port $port"

    echo "Reading all payments"
    curl -s -H "Authorization: Bearer $jwt_token" $uri/payments

    echo "Creating a payment"
    curl -s -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $jwt_token" -d '{"id":"1","order_id":"1","amount":100.0,"status":"pending"}' $uri/payments

    echo "Reading all payments"
    curl -s -H "Authorization: Bearer $jwt_token" $uri/payments

    echo "Updating a payment"
    curl -s -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer $jwt_token" -d '{"id":"1","order_id":"1","amount":150.0,"status":"completed"}' $uri/payments/1

    echo "Reading the updated payment"
    curl -s -H "Authorization: Bearer $jwt_token" $uri/payments/1

    echo "Deleting a payment"
    curl -s -X DELETE -H "Authorization: Bearer $jwt_token" $uri/payments/1

    echo "Reading all payments"
    curl -s -H "Authorization: Bearer $jwt_token" $uri/payments
}

notification_service_tests() {
    local port=8085
    local uri=http://localhost:$port

    echo "Testing notification-service on port $port"

    echo "Reading all notifications"
    curl -s -H "Authorization: Bearer $jwt_token" $uri/notifications

    echo "Creating a notification"
    curl -s -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $jwt_token" -d '{"id":"1","user_id":"1","message":"Your order has been shipped","status":"pending"}' $uri/notifications

    echo "Reading all notifications"
    curl -s -H "Authorization: Bearer $jwt_token" $uri/notifications

    echo "Updating a notification"
    curl -s -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer $jwt_token" -d '{"id":"1","user_id":"1","message":"Your order has been delivered","status":"completed"}' $uri/notifications/1

    echo "Reading the updated notification"
    curl -s -H "Authorization: Bearer $jwt_token" $uri/notifications/1

    echo "Deleting a notification"
    curl -s -X DELETE -H "Authorization: Bearer $jwt_token" $uri/notifications/1

    echo "Reading all notifications"
    curl -s -H "Authorization: Bearer $jwt_token" $uri/notifications
}

# Initialiser la base de données
# initialize_db
# Obtenir un token JWT pour les tests
jwt_token=$(get_jwt_token "user1@example.com" "password")
echo "JWT token: $jwt_token"
# Tester les services
product_service_tests
user_service_tests
order_service_tests
payment_service_tests
notification_service_tests