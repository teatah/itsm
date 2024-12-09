package dashboard

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"html/template"
	"itsm/models"
	_ "itsm/session"
	"itsm/utils"
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
	isClient, err := utils.IsClientUser(r)
	if err != nil {
		http.Error(w, "Ошибка получения сессии: "+err.Error(), http.StatusUnauthorized)
		return
	}

	tmpl, err := template.ParseFiles("templates/dashboard/dashboard.html",
		"templates/header/header.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, map[string]interface{}{
		"IsClient": isClient,
	})
	if err != nil {
		http.Error(w, "Ошибка при выполнении шаблона", http.StatusInternalServerError)
		return
	}
}

func incidentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Раздел Инциденты"))
}

func businessServicesHandler(w http.ResponseWriter, r *http.Request) {
	var services []models.Service

	curSession, err := utils.GetCurSession(r)
	if err != nil {
		http.Error(w, "Ошибка получения сессии: "+err.Error(), http.StatusUnauthorized)
		return
	}

	isClient, err := utils.IsClientUser(r)
	if err != nil {
		http.Error(w, "Ошибка получения сессии: "+err.Error(), http.StatusUnauthorized)
		return
	}

	isAdmin := curSession.Values["isAdmin"].(bool)

	if err := db.Where("is_business = ?", true).Find(&services).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/services/services.html",
		"templates/header/header.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Передаем данные в шаблон, включая права доступа
	err = tmpl.Execute(w, map[string]interface{}{
		"Services":   services,
		"IsAdmin":    isAdmin, // Передаем информацию о правах доступа
		"IsBusiness": true,
		"IsClient":   isClient,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func technicalServicesHandler(w http.ResponseWriter, r *http.Request) {
	var services []models.Service

	curSession, err := utils.GetCurSession(r)
	if err != nil {
		http.Error(w, "Ошибка получения сессии: "+err.Error(), http.StatusUnauthorized)
		return
	}

	isAdmin := curSession.Values["isAdmin"].(bool)

	if err := db.Where("is_technical = ?", true).Find(&services).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/services/services.html",
		"templates/header/header.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Передаем данные в шаблон, включая права доступа
	err = tmpl.Execute(w, map[string]interface{}{
		"Services":    services,
		"IsAdmin":     isAdmin, // Передаем информацию о правах доступа
		"IsTechnical": true,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func messengerHandler(w http.ResponseWriter, r *http.Request) {
	curSession, err := utils.GetCurSession(r)
	if err != nil {
		http.Error(w, "Ошибка получения сессии: "+err.Error(), http.StatusUnauthorized)
		return
	}

	userID, ok := curSession.Values["userID"].(uint)
	if !ok {
		http.Error(w, "User not found in session", http.StatusUnauthorized)
		return
	}

	// Передаем userID в шаблон
	data := struct {
		UserID uint
	}{
		UserID: userID,
	}

	tmpl, err := template.ParseFiles("templates/messenger/messenger.html",
		"templates/header/header.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		return
	}
}
