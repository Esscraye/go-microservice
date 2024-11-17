package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Payment struct {
	ID      uint    `gorm:"primaryKey;autoIncrement" json:"id"`
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
		getPayments(w)
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
		getPayment(w, id)
	case "PUT":
		updatePayment(w, r, id)
	case "DELETE":
		deletePayment(w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getPayments(w http.ResponseWriter) {
	var payments []Payment
	if err := db.Find(&payments).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(payments)
}

func createPayment(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var payment Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Vérifier l'existence de la commande
	orderExists := checkOrderExists(payment.OrderID, token)
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

func checkOrderExists(orderID string, token string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://order-service:8083/orders/%s", orderID), nil)
	if err != nil {
		return false
	}
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func getPayment(w http.ResponseWriter, id string) {
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

func deletePayment(w http.ResponseWriter, id string) {
	if err := db.Delete(&Payment{}, "id = ?", id).Error; err != nil {
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
	mux.Handle("/payments", jwtMiddleware(http.HandlerFunc(paymentsHandler)))
	mux.Handle("/payments/", jwtMiddleware(http.HandlerFunc(paymentHandler)))
	mux.HandleFunc("/health", healthHandler)

	log.Println("Starting Payment Service on :8084")
	log.Fatal(http.ListenAndServe(":8084", mux))
}
