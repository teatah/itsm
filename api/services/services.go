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
	r.HandleFunc("/delete-service", deleteServiceHandler).Methods("DELETE")
	r.HandleFunc("/edit-service", editServiceHandler).Methods("GET")
	r.HandleFunc("/open-service", openServiceHandler).Methods("GET")
	r.HandleFunc("/add-service", addServiceHandler).Methods("GET")
	r.HandleFunc("/update-service", updateServiceHandler).Methods("PUT")
	r.HandleFunc("/create-service", createServiceHandler).Methods("POST")
}

func updateServiceHandler(w http.ResponseWriter, r *http.Request) {
	serviceID := r.URL.Query().Get("id")
	if serviceID == "" {
		http.Error(w, "ID услуги не указан", http.StatusBadRequest)
		return
	}

	var service models.Service
	if err := db.First(&service, serviceID).Error; err != nil {
		http.Error(w, "Услуга не найдена", http.StatusNotFound)
		return
	}

	// Обновляем поля услуги
	service.Name = r.FormValue("name")
	service.Description = r.FormValue("description")
	service.IsBusiness = r.FormValue("serviceType") == "business"
	service.IsTechnical = r.FormValue("serviceType") == "technical"

	// Сохраняем изменения в базе данных
	if err := db.Save(&service).Error; err != nil {
		http.Error(w, "Ошибка при обновлении услуги", http.StatusInternalServerError)
		return
	}

	// Перенаправляем на список услуг
	http.Redirect(w, r, "/business-services", http.StatusSeeOther)
}

func createServiceHandler(w http.ResponseWriter, r *http.Request) {
	// Создаем новую услугу
	service := models.Service{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		IsBusiness:  r.FormValue("serviceType") == "business",
		IsTechnical: r.FormValue("serviceType") == "technical",
	}

	// Сохраняем новую услугу в базе данных
	if err := db.Create(&service).Error; err != nil {
		http.Error(w, "Ошибка при создании услуги", http.StatusInternalServerError)
		return
	}

	// Перенаправляем на список услуг
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
	serviceID := r.URL.Query().Get("id")
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
	serviceID := r.URL.Query().Get("id")
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
	// Извлечение ID услуги из параметров запроса
	serviceID := r.URL.Query().Get("id")
	if serviceID == "" {
		http.Error(w, "ID услуги не указан", http.StatusBadRequest)
		return
	}

	// Преобразование ID в тип uint
	id, err := strconv.ParseUint(serviceID, 10, 32)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	// Удаление услуги из базы данных
	var service models.Service
	if err := db.First(&service, id).Error; err != nil {
		http.Error(w, "Услуга не найдена", http.StatusNotFound)
		return
	}

	if err := db.Delete(&service).Error; err != nil {
		http.Error(w, "Ошибка при удалении услуги", http.StatusInternalServerError)
		return
	}

	// Успешное удаление
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
