package services

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"html/template"
	"itsm/models"
	"itsm/utils"
	"net/http"
	"strconv"
)

var db *gorm.DB

func SetupRoutes(r *mux.Router, database *gorm.DB) {
	db = database
	r.HandleFunc("/service/{id}/delete", deleteServiceHandler).Methods("DELETE")
	r.HandleFunc("/service/{id}/edit", editServiceHandler).Methods("GET")
	r.HandleFunc("/service/{id}", openServiceHandler).Methods("GET")
	r.HandleFunc("/service/{id}/update", updateServiceHandler).Methods("PUT")
	r.HandleFunc("/services/create", createServiceHandler).Methods("POST")
	r.HandleFunc("/services/add", addServiceHandler).Methods("GET")
}

func updateServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceID := vars["id"]

	if serviceID == "" {
		http.Error(w, "ID услуги не указан", http.StatusBadRequest)
		return
	}

	var service models.Service
	if err := db.First(&service, serviceID).Error; err != nil {
		http.Error(w, "Услуга не найдена", http.StatusNotFound)
		return
	}

	service.Name = r.FormValue("name")
	service.Description = r.FormValue("description")
	service.IsBusiness = r.FormValue("serviceType") == "business"
	service.IsTechnical = r.FormValue("serviceType") == "technical"

	if err := db.Save(&service).Error; err != nil {
		http.Error(w, "Ошибка при обновлении услуги", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/business-services", http.StatusSeeOther)
}

func createServiceHandler(w http.ResponseWriter, r *http.Request) {
	service := models.Service{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		IsBusiness:  r.FormValue("serviceType") == "business",
		IsTechnical: r.FormValue("serviceType") == "technical",
	}

	if err := db.Create(&service).Error; err != nil {
		http.Error(w, "Ошибка при создании услуги", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/business-services", http.StatusSeeOther)
}

func addServiceHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Service  models.Service
		IsCreate bool
		IsView   bool
		IsEdit   bool
		IsClient bool
	}{
		Service:  models.Service{},
		IsCreate: true,
		IsView:   false,
		IsEdit:   false,
		IsClient: false,
	}

	renderTemplate(w, "templates/service/service.html", data)
}

func openServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceID := vars["id"]

	if serviceID == "" {
		http.Error(w, "ID услуги не указан", http.StatusBadRequest)
		return
	}

	var service models.Service
	if err := db.First(&service, serviceID).Error; err != nil {
		http.Error(w, "Услуга не найдена", http.StatusNotFound)
		return
	}

	isClient, err := utils.IsClientUser(r)
	if err != nil {
		http.Error(w, "Ошибка получения сессии: "+err.Error(), http.StatusUnauthorized)
		return
	}

	data := struct {
		Service  models.Service
		IsCreate bool
		IsView   bool
		IsEdit   bool
		IsClient bool
	}{
		Service:  service,
		IsCreate: false,
		IsView:   true,
		IsEdit:   false,
		IsClient: isClient,
	}

	renderTemplate(w, "templates/service/service.html", data)
}

func editServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceID := vars["id"]

	if serviceID == "" {
		http.Error(w, "ID услуги не указан", http.StatusBadRequest)
		return
	}

	var service models.Service
	if err := db.First(&service, serviceID).Error; err != nil {
		http.Error(w, "Услуга не найдена", http.StatusNotFound)
		return
	}

	data := struct {
		Service  models.Service
		IsCreate bool
		IsView   bool
		IsEdit   bool
		IsClient bool
	}{
		Service:  service,
		IsCreate: false,
		IsView:   false,
		IsEdit:   true,
		IsClient: false,
	}

	renderTemplate(w, "templates/service/service.html", data)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl, "templates/header/header.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func deleteServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceID := vars["id"]

	if serviceID == "" {
		http.Error(w, "ID услуги не указан", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(serviceID, 10, 32)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	var service models.Service
	if err := db.First(&service, id).Error; err != nil {
		http.Error(w, "Услуга не найдена", http.StatusNotFound)
		return
	}

	if err := db.Delete(&service).Error; err != nil {
		http.Error(w, "Ошибка при удалении услуги", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
