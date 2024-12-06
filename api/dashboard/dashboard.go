package dashboard

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"html/template"
	"itsm/models"
	"net/http"
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
	err = tmpl.Execute(w, nil)
	if err != nil {
		return
	}
}

func businessServicesHandler(w http.ResponseWriter, r *http.Request) {
	var services []models.Service
	if err := db.Preload("ServiceLine").
		Joins("JOIN service_lines ON service_lines.id = services.service_line_id").
		Where("service_lines.name = ?", "Бизнес услуги").
		Find(&services).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/business_services/business_services.html",
		"templates/header/header.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, services)
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
