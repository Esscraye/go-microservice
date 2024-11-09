package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"auth-service/pkg"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID       string `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var (
	db     *gorm.DB
	jwtKey = []byte(os.Getenv("JWT_SECRET"))
)

func initDB() {
	dsn := "host=db user=user password=password dbname=microservices port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
}

func migrate() {
	db.AutoMigrate(&User{})
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getUsers(w, r)
	case "POST":
		createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	switch r.Method {
	case "GET":
		getUser(w, r, id)
	case "PUT":
		updateUser(w, r, id)
	case "DELETE":
		deleteUser(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := db.Create(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func getUser(w http.ResponseWriter, r *http.Request, id string) {
	var user User
	if err := db.First(&user, "id = ?", id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request, id string) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := db.Model(&User{}).Where("id = ?", id).Updates(user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func deleteUser(w http.ResponseWriter, r *http.Request, id string) {
	if err := db.Delete(&User{}, "id = ?", id).Error; err != nil {
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

	http.Handle("/users", pkg.AuthMiddleware(http.HandlerFunc(usersHandler)))
	http.Handle("/users/", pkg.AuthMiddleware(http.HandlerFunc(userHandler)))
	http.HandleFunc("/health", healthHandler)

	log.Println("Starting User Service on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
