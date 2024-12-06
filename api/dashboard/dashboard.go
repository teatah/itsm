package dashboard

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"html/template"
	"itsm/models"
	"itsm/session"
	"itsm/utils"
	"net/http"

	_ "itsm/session"
)

var db *gorm.DB

func SetupRoutes(r *mux.Router, database *gorm.DB) {
	db = database
	r.HandleFunc("/dashboard", dashboardHandler)
	r.HandleFunc("/business-services", businessServicesHandler)
	r.HandleFunc("/technical-services", technicalServicesHandler)
	r.HandleFunc("/incidents", incidentsHandler)
	r.HandleFunc("/messenger", messengerHandler)
	r.HandleFunc("/users/get", getUsersHandler)
}

func getUsersHandler(w http.ResponseWriter, _ *http.Request) {
	var users []models.User

	// Получаем всех пользователей
	if err := db.Find(&users).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем ID пользователей, с которыми уже есть сообщения
	var userIDsWithMessages []uint
	if err := db.Model(&models.Message{}).Select("DISTINCT sender_id").Scan(&userIDsWithMessages).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Фильтруем пользователей, исключая тех, у кого есть сообщения
	var filteredUsers []models.User
	for _, user := range users {
		if !contains(userIDsWithMessages, user.ID) {
			filteredUsers = append(filteredUsers, user)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(filteredUsers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func contains(slice []uint, item uint) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func dashboardHandler(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("templates/dashboard/dashboard.html",
		"templates/header/header.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		return
	}
}

func businessServicesHandler(w http.ResponseWriter, r *http.Request) {
	var services []models.Service

	port := utils.GetPort(r)
	sessionName := "session-" + port
	curSession, err := session.Store.Get(r, sessionName)

	isAdmin := curSession.Values["isAdmin"].(bool)

	if err := db.Where("is_business = ?", true).Find(&services).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/business_services/business_services.html",
		"templates/header/header.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Передаем данные в шаблон, включая права доступа
	err = tmpl.Execute(w, map[string]interface{}{
		"Services": services,
		"IsAdmin":  isAdmin, // Передаем информацию о правах доступа
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func technicalServicesHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Раздел Технические услуги"))
}

func incidentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Раздел Инциденты"))
}

func messengerHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/messenger/messenger.html",
		"templates/header/header.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		return
	}
}
