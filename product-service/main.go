package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Product struct {
	ID       string `gorm:"primaryKey"`
	Name     string
	Category string
	Price    float64
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
	db.AutoMigrate(&Product{})
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getProducts(w)
	case "POST":
		createProduct(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/products/"):]
	switch r.Method {
	case "GET":
		getProduct(w, id)
	case "PUT":
		updateProduct(w, r, id)
	case "DELETE":
		deleteProduct(w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getProducts(w http.ResponseWriter) {
	var products []Product
	if err := db.Find(&products).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(products)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := db.Create(&product).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func getProduct(w http.ResponseWriter, id string) {
	var product Product
	if err := db.First(&product, "id = ?", id).Error; err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(product)
}

func updateProduct(w http.ResponseWriter, r *http.Request, id string) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := db.Model(&Product{}).Where("id = ?", id).Updates(product).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func deleteProduct(w http.ResponseWriter, id string) {
	if err := db.Delete(&Product{}, "id = ?", id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type contextKey string

const userIDKey contextKey = "user_id"

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
		ctx := context.WithValue(r.Context(), userIDKey, result.UserID)
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
	mux.Handle("/products", jwtMiddleware(http.HandlerFunc(productsHandler)))
	mux.Handle("/products/", jwtMiddleware(http.HandlerFunc(productHandler)))
	mux.HandleFunc("/health", healthHandler)

	log.Println("Starting Product Service on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
