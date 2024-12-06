package dashboard

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"html/template"
	"itsm/models"
	"itsm/session"
	"net/http"
	"strings"

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
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/dashboard/dashboard.html",
		"templates/header/header.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func businessServicesHandler(w http.ResponseWriter, r *http.Request) {
	var services []models.Service

	// Получаем порт из r.Host
	hostParts := strings.Split(r.Host, ":")
	var port string
	if len(hostParts) > 1 {
		port = hostParts[1] // Порт будет вторым элементом
	} else {
		// Если порт не указан, устанавливаем значение по умолчанию
		if r.URL.Scheme == "http" {
			port = "80"
		} else if r.URL.Scheme == "https" {
			port = "443"
		}
	}

	sessionName := "session-" + port

	session, err := session.Store.Get(r, sessionName)

	isAdmin := session.Values["isAdmin"].(bool)

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
	w.Write([]byte("Раздел Мессенджер"))
}
