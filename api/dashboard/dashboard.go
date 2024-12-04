package dashboard

import (
	"database/sql"
	"html/template"
	"net/http"
)

var db *sql.DB

func SetupRoutes(database *sql.DB) {
	db = database
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/business-services", businessServicesHandler)
	http.HandleFunc("/technical-services", technicalServicesHandler)
	http.HandleFunc("/incidents", incidentsHandler)
	http.HandleFunc("/messenger", messengerHandler)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/dashboard.html"))
	tmpl.Execute(w, nil)
}

func businessServicesHandler(w http.ResponseWriter, r *http.Request) {
	// Здесь будет логика для отображения бизнес услуг
	w.Write([]byte("Раздел Бизнес услуги"))
}

func technicalServicesHandler(w http.ResponseWriter, r *http.Request) {
	// Здесь будет логика для отображения технических услуг
	w.Write([]byte("Раздел Технические услуги"))
}

func incidentsHandler(w http.ResponseWriter, r *http.Request) {
	// Здесь будет логика для отображения инцидентов
	w.Write([]byte("Раздел Инциденты"))
}

func messengerHandler(w http.ResponseWriter, r *http.Request) {
	// Здесь будет логика для отображения мессенджера
	w.Write([]byte("Раздел Мессенджер"))
}
