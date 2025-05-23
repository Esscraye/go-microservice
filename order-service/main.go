package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	productServiceURL = os.Getenv("PRODUCT_SERVICE_URL")
	jwtKey            = []byte(os.Getenv("JWT_SECRET"))
)

type Order struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Status    string `json:"status"`
}

var db *gorm.DB

func initDB() {
	dsn := "host=db user=user password=password dbname=microservices port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
}

func migrate() {
	db.AutoMigrate(&Order{})
}

func ordersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getOrders(w)
	case "POST":
		createOrder(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/orders/"):]
	switch r.Method {
	case "GET":
		getOrder(w, id)
	case "PUT":
		updateOrder(w, r, id)
	case "DELETE":
		deleteOrder(w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getOrders(w http.ResponseWriter) {
	var orders []Order
	if err := db.Find(&orders).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(orders)
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Vérifier la disponibilité du produit
	productAvailable := checkProductAvailability(order.ProductID, token)
	if !productAvailable {
		http.Error(w, "Product not available", http.StatusBadRequest)
		return
	}

	if err := db.Create(&order).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func checkProductAvailability(productID string, token string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	for i := 0; i < 3; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/products/%s", productServiceURL, productID), nil)
		if err != nil {
			continue
		}
		req.Header.Set("Authorization", token)

		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			return true
		}
		time.Sleep(2 * time.Second)
	}
	return false
}

func getOrder(w http.ResponseWriter, id string) {
	var order Order
	if err := db.First(&order, "id = ?", id).Error; err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(order)
}

func updateOrder(w http.ResponseWriter, r *http.Request, id string) {
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := db.Model(&Order{}).Where("id = ?", id).Updates(order).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func deleteOrder(w http.ResponseWriter, id string) {
	if err := db.Delete(&Order{}, "id = ?", id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		resp, err := http.Post("http://auth-service:8080/verify-token", "application/json", strings.NewReader(fmt.Sprintf(`{"token":"%s"}`, token)))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var result struct {
			UserID string `json:"user_id"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user_id", result.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	initDB()
	migrate()

	mux := http.NewServeMux()
	mux.Handle("/orders", jwtMiddleware(http.HandlerFunc(ordersHandler)))
	mux.Handle("/orders/", jwtMiddleware(http.HandlerFunc(orderHandler)))
	mux.HandleFunc("/health", healthHandler)

	log.Println("Starting Order Service on :8083")
	log.Fatal(http.ListenAndServe(":8083", mux))
}
