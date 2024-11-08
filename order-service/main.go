package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	productServiceURL = os.Getenv("PRODUCT_SERVICE_URL")
	jwtKey            = []byte(os.Getenv("JWT_SECRET"))
)

type Order struct {
	ID        string `gorm:"primaryKey" json:"id"`
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
		getOrders(w, r)
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
		getOrder(w, r, id)
	case "PUT":
		updateOrder(w, r, id)
	case "DELETE":
		deleteOrder(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	var orders []Order
	if err := db.Find(&orders).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(orders)
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Vérifier la disponibilité du produit
	productAvailable := checkProductAvailability(order.ProductID)
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

func checkProductAvailability(productID string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	for i := 0; i < 3; i++ {
		resp, err := client.Get(fmt.Sprintf("%s/products/%s", productServiceURL, productID))
		if err == nil && resp.StatusCode == http.StatusOK {
			return true
		}
		time.Sleep(2 * time.Second)
	}
	return false
}

func getOrder(w http.ResponseWriter, r *http.Request, id string) {
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

func deleteOrder(w http.ResponseWriter, r *http.Request, id string) {
	if err := db.Delete(&Order{}, "id = ?", id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	initDB()
	migrate()

	http.HandleFunc("/orders", ordersHandler)
	http.HandleFunc("/orders/", orderHandler)
	http.HandleFunc("/health", healthHandler)

	log.Println("Starting Order Service on :8083")
	log.Fatal(http.ListenAndServe(":8083", nil))
}
