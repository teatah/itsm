package incidents

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
	r.HandleFunc("/incidents/add", addIncidentHandler).Methods("GET")
	r.HandleFunc("/incidents/add", createIncidentHandler).Methods("POST")
	r.HandleFunc("/incident/{id}", incidentHandler).Methods("GET")
	r.HandleFunc("/incident/{id}/update", incidentHandler).Methods("POST")
}

func addIncidentHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/incidents/incident_add/add_incident.html",
		"templates/header/header.html")
	if err != nil {
		http.Error(w, "Ошибка при загрузке шаблона", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Ошибка при выполнении шаблона", http.StatusInternalServerError)
	}
}

func createIncidentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	status := r.FormValue("status")

	curSession, err := utils.GetCurSession(r)
	if err != nil {
		http.Error(w, "Ошибка при получении сессии", http.StatusInternalServerError)
		return
	}

	userID, ok := curSession.Values["userID"].(uint)
	if !ok {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	// Создаем новый инцидент
	incident := models.Incident{
		Title:       title,
		Description: description,
		Status:      status,
		UserID:      userID,
	}

	// Сохраняем инцидент в базе данных
	if err := db.Create(&incident).Error; err != nil {
		http.Error(w, "Ошибка при добавлении инцидента", http.StatusInternalServerError)
		return
	}

	// Перенаправляем на страницу со списком инцидентов
	http.Redirect(w, r, "/incidents", http.StatusSeeOther)
}

func incidentHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID инцидента из URL
	vars := mux.Vars(r) // Получаем переменные из URL
	id := vars["id"]    // Извлекаем ID инцидента

	// Проверяем, что ID не пустой
	if id == "" {
		http.Error(w, "ID инцидента не указан", http.StatusBadRequest)
		return
	}

	var incident models.Incident
	var user models.User

	if err := db.First(&incident, id).Error; err != nil {
		http.Error(w, "Инцидент не найден", http.StatusNotFound)
		return
	}

	if err := db.First(&user, incident.UserID).Error; err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	curSession, err := utils.GetCurSession(r)
	if err != nil {
		http.Error(w, "Ошибка получения сессии: "+err.Error(), http.StatusUnauthorized)
		return
	}
	isAdmin := curSession.Values["isAdmin"].(bool)
	isTechOfficer := curSession.Values["isTechOfficer"].(bool)

	// Обрабатываем данные из формы, если метод POST
	if r.Method == http.MethodPost {
		status := r.FormValue("status")
		responsibleUserID := r.FormValue("responsible_user_id")

		// Обновляем статус инцидента
		incident.Status = status

		// Преобразуем ID ответственного пользователя в uint
		if responsibleUserID != "" {
			var responsibleUser models.User
			if err := db.First(&responsibleUser, responsibleUserID).Error; err == nil {
				incident.ResponsibleUserID = responsibleUser.ID
			}
		}

		// Сохраняем изменения в базе данных
		if err := db.Save(&incident).Error; err != nil {
			http.Error(w, "Ошибка при обновлении инцидента", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/incidents", http.StatusSeeOther)
		return
	}

	// Загружаем список пользователей с ролью тех.поддержки, если пользователь имеет соответствующие права
	var techOfficers []models.User
	if isAdmin || isTechOfficer {
		if err := db.Where("is_tech_officer = ?", true).Find(&techOfficers).Error; err != nil {
			http.Error(w, "Ошибка при загрузке пользователей", http.StatusInternalServerError)
			return
		}
	}

	// Создаем данные для передачи в шаблон
	data := map[string]interface{}{
		"Incident":     incident,
		"Username":     user.Username,
		"TechOfficers": techOfficers,
		"IsEditable":   isAdmin || isTechOfficer, // Флаг для отображения формы редактирования
	}

	// Загружаем и выполняем шаблон
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
