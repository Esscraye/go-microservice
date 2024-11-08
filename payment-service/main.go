package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Payment struct {
	ID      string  `gorm:"primaryKey" json:"id"`
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
	Status  string  `json:"status"`
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
	db.AutoMigrate(&Payment{})
}

func paymentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getPayments(w, r)
	case "POST":
		createPayment(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/payments/"):]
	switch r.Method {
	case "GET":
		getPayment(w, r, id)
	case "PUT":
		updatePayment(w, r, id)
	case "DELETE":
		deletePayment(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getPayments(w http.ResponseWriter, r *http.Request) {
	var payments []Payment
	if err := db.Find(&payments).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(payments)
}

func createPayment(w http.ResponseWriter, r *http.Request) {
	var payment Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Vérifier l'existence de la commande
	orderExists := checkOrderExists(payment.OrderID)
	if !orderExists {
		http.Error(w, "Order not found", http.StatusBadRequest)
		return
	}

	if err := db.Create(&payment).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func checkOrderExists(orderID string) bool {
	resp, err := http.Get(fmt.Sprintf("http://order-service:8083/orders/%s", orderID))
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func getPayment(w http.ResponseWriter, r *http.Request, id string) {
	var payment Payment
	if err := db.First(&payment, "id = ?", id).Error; err != nil {
		http.Error(w, "Payment not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(payment)
}

func updatePayment(w http.ResponseWriter, r *http.Request, id string) {
	var payment Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := db.Model(&Payment{}).Where("id = ?", id).Updates(payment).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func deletePayment(w http.ResponseWriter, r *http.Request, id string) {
	if err := db.Delete(&Payment{}, "id = ?", id).Error; err != nil {
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

	http.HandleFunc("/payments", paymentsHandler)
	http.HandleFunc("/payments/", paymentHandler)
	http.HandleFunc("/health", healthHandler)

	log.Println("Starting Payment Service on :8084")
	log.Fatal(http.ListenAndServe(":8084", nil))
}
