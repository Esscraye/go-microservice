package main

import (
	"encoding/json"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Notification struct {
	ID      string `gorm:"primaryKey" json:"id"`
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	Status  string `json:"status"`
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
	db.AutoMigrate(&Notification{})
}

func notificationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getNotifications(w, r)
	case "POST":
		createNotification(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func notificationHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/notifications/"):]
	switch r.Method {
	case "GET":
		getNotification(w, r, id)
	case "PUT":
		updateNotification(w, r, id)
	case "DELETE":
		deleteNotification(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getNotifications(w http.ResponseWriter, r *http.Request) {
	var notifications []Notification
	if err := db.Find(&notifications).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(notifications)
}

func createNotification(w http.ResponseWriter, r *http.Request) {
	var notification Notification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := db.Create(&notification).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	sendNotification(notification)
}

func getNotification(w http.ResponseWriter, r *http.Request, id string) {
	var notification Notification
	if err := db.First(&notification, "id = ?", id).Error; err != nil {
		http.Error(w, "Notification not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(notification)
}

func updateNotification(w http.ResponseWriter, r *http.Request, id string) {
	var notification Notification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := db.Model(&Notification{}).Where("id = ?", id).Updates(notification).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func deleteNotification(w http.ResponseWriter, r *http.Request, id string) {
	if err := db.Delete(&Notification{}, "id = ?", id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func sendNotification(notification Notification) {
	// Logic to send notification via email or SMS
	log.Printf("Sending notification to user %s: %s\n", notification.UserID, notification.Message)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	initDB()
	migrate()

	http.HandleFunc("/notifications", notificationsHandler)
	http.HandleFunc("/notifications/", notificationHandler)
	http.HandleFunc("/health", healthHandler)

	log.Println("Starting Notification Service on :8085")
	log.Fatal(http.ListenAndServe(":8085", nil))
}
