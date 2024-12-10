package incidents

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"html/template"
	"itsm/models"
	"itsm/utils"
	"net/http"
	"strconv"
	"strings"
)

var db *gorm.DB

func SetupRoutes(r *mux.Router, database *gorm.DB) {
	db = database
	r.HandleFunc("/incidents/add", addIncidentHandler).Methods("GET")
	r.HandleFunc("/incidents/create", createIncidentHandler).Methods("POST")
	r.HandleFunc("/incident/{id}", incidentHandler).Methods("GET")
	r.HandleFunc("/incident/{id}/update", updateIncidentsHandler).Methods("POST")
}

func addIncidentHandler(w http.ResponseWriter, r *http.Request) {
	isClient, err := utils.IsClientUser(r)
	if err != nil {
		http.Error(w, "Ошибка получения сессии: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Получаем список услуг из базы данных
	var services []models.Service
	if err := db.Where("is_business = ?", true).Find(&services).Error; err != nil {
		http.Error(w, "Ошибка при получении услуг: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		IsClient bool
		Services []models.Service
	}{
		IsClient: isClient,
		Services: services,
	}

	tmpl, err := template.ParseFiles("templates/incidents/incident_add/add_incident.html",
		"templates/header/header.html")
	if err != nil {
		http.Error(w, "Ошибка при загрузке шаблона", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Ошибка при выполнении шаблона", http.StatusInternalServerError)
	}
}

func createIncidentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Получаем текущую сессию
	curSession, err := utils.GetCurSession(r)
	if err != nil {
		http.Error(w, "Ошибка при получении сессии", http.StatusInternalServerError)
		return
	}

	// Получаем текущего пользователя
	userID, ok := curSession.Values["userID"].(uint)
	if !ok {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	// Создаем новый инцидент
	incident := models.Incident{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Status:      "Открыт",
		UserID:      userID,
	}

	// Сохраняем инцидент в базе данных
	if err := db.Create(&incident).Error; err != nil {
		http.Error(w, "Ошибка при добавлении инцидента", http.StatusInternalServerError)
		return
	}

	// Получаем выбранные услуги из формы
	serviceIDs := r.FormValue("selected_services") // Получаем строку с ID услуг

	if serviceIDs != "" {
		ids := strings.Split(serviceIDs, ",") // Разделяем строку на массив ID

		// Привязываем выбранные услуги к инциденту
		for _, idStr := range ids {
			id, err := strconv.ParseUint(idStr, 10, 32) // Преобразуем строку в uint
			if err != nil {
				http.Error(w, "Ошибка при преобразовании ID услуги", http.StatusBadRequest)
				return
			}

			incidentService := models.IncidentService{
				IncidentID: incident.ID,
				ServiceID:  uint(id),
			}

			// Сохраняем связь в базе данных
			if err := db.Create(&incidentService).Error; err != nil {
				http.Error(w, "Ошибка при добавлении связи инцидента и услуги", http.StatusInternalServerError)
				return
			}
		}
	}

	// Перенаправляем на страницу со списком инцидентов
	http.Redirect(w, r, "/incidents", http.StatusSeeOther)
}

func incidentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "ID инцидента не указан", http.StatusBadRequest)
		return
	}

	var incident models.Incident
	if err := db.First(&incident, id).Error; err != nil {
		http.Error(w, "Инцидент не найден", http.StatusNotFound)
		return
	}

	var user models.User
	if err := db.First(&user, incident.UserID).Error; err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	var responsibleUserUsername string

	if incident.ResponsibleUserID == nil {
		responsibleUserUsername = "Не назначен"
	} else {
		var responsibleUser models.User
		if err := db.First(&responsibleUser, *incident.ResponsibleUserID).Error; err != nil {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}
		responsibleUserUsername = responsibleUser.Username
	}

	curSession, err := utils.GetCurSession(r)
	if err != nil {
		http.Error(w, "Ошибка получения сессии: "+err.Error(), http.StatusUnauthorized)
		return
	}
	isAdmin := curSession.Values["isAdmin"].(bool)
	isTechOfficer := curSession.Values["isTechOfficer"].(bool)
	isClient, err := utils.IsClientUser(r)
	if err != nil {
		http.Error(w, "Ошибка получения сессии: "+err.Error(), http.StatusUnauthorized)
		return
	}

	var techOfficers []models.User
	if isAdmin || isTechOfficer {
		if err := db.Where("is_tech_officer = ?", true).Find(&techOfficers).Error; err != nil {
			http.Error(w, "Ошибка при загрузке пользователей", http.StatusInternalServerError)
			return
		}
	}

	var services []models.Service
	if err := db.Where("is_business = ?", true).Find(&services).Error; err != nil {
		http.Error(w, "Ошибка при получении услуг: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем только ID выбранных услуг для данного инцидента
	var selectedServiceIDs []uint
	if err := db.Model(&models.IncidentService{}).
		Select("service_id").
		Where("incident_id = ?", incident.ID).
		Find(&selectedServiceIDs).Error; err != nil {
		http.Error(w, "Ошибка при получении услуг: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаем новый срез для хранения только выбранных услуг
	var selectedServices []models.Service
	for _, service := range services {
		for _, selectedID := range selectedServiceIDs {
			if service.ID == selectedID {
				selectedServices = append(selectedServices, service)
				break // Выходим из внутреннего цикла, если нашли совпадение
			}
		}
	}

	var responsibleUserIDValue uint
	if incident.ResponsibleUserID != nil {
		responsibleUserIDValue = *incident.ResponsibleUserID
	}

	data := map[string]interface{}{
		"Incident":                incident,
		"Username":                user.Username,
		"ResponsibleUserID":       responsibleUserIDValue,
		"ResponsibleUserUsername": responsibleUserUsername,
		"TechOfficers":            techOfficers,
		"Services":                services,
		"SelectedServices":        selectedServices,
		"HasEditRights":           isAdmin || isTechOfficer,
		"IsClient":                isClient,
	}

	tmpl, err := template.ParseFiles("templates/incidents/incident/incident.html",
		"templates/header/header.html")
	if err != nil {
		http.Error(w, "Ошибка при загрузке шаблона", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Ошибка при выполнении шаблона", http.StatusInternalServerError)
		return
	}
}

func updateIncidentsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "ID инцидента не указан", http.StatusBadRequest)
		return
	}

	var incident models.Incident
	if err := db.First(&incident, id).Error; err != nil {
		http.Error(w, "Инцидент не найден", http.StatusNotFound)
		return
	}

	status := r.FormValue("status")
	responsibleUserID := r.FormValue("responsible_user_id")

	incident.Status = status
	if responsibleUserID != "" {
		userID, err := strconv.ParseUint(responsibleUserID, 10, 32)
		if err == nil {
			incident.ResponsibleUserID = new(uint)
			*incident.ResponsibleUserID = uint(userID)
		}
	} else {
		incident.ResponsibleUserID = nil
	}

	// Удаляем старые связи
	if err := db.Where("incident_id = ?", incident.ID).Delete(&models.IncidentService{}).Error; err != nil {
		http.Error(w, "Ошибка при удалении старых услуг", http.StatusInternalServerError)
		return
	}

	// Обработка выбранных услуг
	selectedServices := r.FormValue("selected_services")
	if selectedServices != "" {
		// Разделяем строку на массив ID услуг
		serviceIDs := strings.Split(selectedServices, ",")

		// Добавляем новые связи
		for _, serviceID := range serviceIDs {
			if serviceID != "" {
				id, err := strconv.ParseUint(serviceID, 10, 32)
				if err == nil {
					incidentService := models.IncidentService{
						IncidentID: incident.ID,
						ServiceID:  uint(id),
					}
					if err := db.Create(&incidentService).Error; err != nil {
						http.Error(w, "Ошибка при добавлении услуги", http.StatusInternalServerError)
						return
					}
				}
			}
		}
	}

	if err := db.Save(&incident).Error; err != nil {
		http.Error(w, "Ошибка при обновлении инцидента", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/incidents", http.StatusSeeOther)
}
