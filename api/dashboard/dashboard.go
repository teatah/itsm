package dashboard

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"html/template"
	"itsm/models"
	"itsm/session"
	"itsm/utils"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "itsm/session"
)

var db *gorm.DB

func SetupRoutes(r *mux.Router, database *gorm.DB) {
	db = database
	r.HandleFunc("/dashboard", dashboardHandler)
	r.HandleFunc("/incidents", incidentsHandler)
	r.HandleFunc("/messenger", messengerHandler)
	r.HandleFunc("/messenger/create", createConversationHandler)
	r.HandleFunc("/conversation", conversationHandler)
	r.HandleFunc("/conversation/send", sendMessageHandler) // Отправка сообщения
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

func incidentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Раздел Инциденты"))
}

/*func messengerHandler(w http.ResponseWriter, r *http.Request) {
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
}*/

func messengerHandler(w http.ResponseWriter, r *http.Request) {
	port := utils.GetPort(r)
	sessionName := "session-" + port
	curSession, _ := session.Store.Get(r, sessionName)

	userID, ok := curSession.Values["userID"].(uint)
	if !ok {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	var conversations []models.Conversation
	// Здесь нужно будет извлечь переписки из базы данных
	if err := db.Where("user1_id = ? OR user2_id = ?", userID, userID).Find(&conversations).Error; err != nil {
		log.Println("Ошибка при извлечении переписок:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Получаем всех пользователей, кроме текущего
	var users []models.User
	if err := db.Where("id != ?", userID).Find(&users).Error; err != nil {
		log.Println("Ошибка при извлечении пользователей:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Получаем текущего пользователя
	var currentUser models.User
	if err := db.First(&currentUser, userID).Error; err != nil {
		log.Println("Ошибка при извлечении текущего пользователя:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Создаем мапу для быстрого доступа к пользователям по ID
	userMap := make(map[uint]string)
	for _, user := range users {
		userMap[user.ID] = user.Username
	}

	tmpl := template.Must(template.ParseFiles("templates/messenger/messenger.html",
		"templates/header/header.html"))
	err := tmpl.Execute(w, map[string]interface{}{
		"Conversations": conversations,
		"Users":         users,
		"UserMap":       userMap,
		"CurrentUser":   currentUser,
	})

	if err != nil {
		log.Println("Ошибка при выполнении шаблона:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}

}

func conversationHandler(w http.ResponseWriter, r *http.Request) {
	// Извлечение ID переписки из URL
	conversationID := r.URL.Query().Get("id")

	var conversation models.Conversation
	if err := db.Preload("Messages").First(&conversation, conversationID).Error; err != nil {
		log.Println("Ошибка при извлечении переписки:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Получаем текущего пользователя
	port := utils.GetPort(r)
	sessionName := "session-" + port
	curSession, _ := session.Store.Get(r, sessionName)
	userID, ok := curSession.Values["userID"].(uint)
	if !ok {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/conversation/conversation.html", "templates/header/header.html"))
	err := tmpl.Execute(w, map[string]interface{}{
		"Conversation": conversation,
		"CurrentUser":  userID,
	})

	if err != nil {
		log.Println("Ошибка при выполнении шаблона:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}

func createConversationHandler(w http.ResponseWriter, r *http.Request) {
	port := utils.GetPort(r)
	sessionName := "session-" + port
	curSession, _ := session.Store.Get(r, sessionName)

	userID, ok := curSession.Values["userID"].(uint)
	if !ok {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		receiverID, err := strconv.ParseUint(r.FormValue("receiver_id"), 10, 32)
		if err != nil {
			http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
			return
		}

		// Проверяем, существует ли уже переписка с этим пользователем
		var existingConversation models.Conversation
		if err := db.Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)", userID, receiverID, receiverID, userID).First(&existingConversation).Error; err == nil {
			// Переписка уже существует, перенаправляем на страницу переписки
			http.Redirect(w, r, fmt.Sprintf("/conversation?id=%d", existingConversation.ID), http.StatusSeeOther)
			return
		}

		// Создаем новую переписку
		newConversation := models.Conversation{
			User1ID: userID,
			User2ID: uint(receiverID),
		}

		if err := db.Create(&newConversation).Error; err != nil {
			log.Println("Ошибка при создании переписки:", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		// Перенаправляем на страницу новой переписки
		http.Redirect(w, r, fmt.Sprintf("/conversation?id=%d", newConversation.ID), http.StatusSeeOther)
		return
	}

	// Если метод не POST, перенаправляем на мессенджер
	http.Redirect(w, r, "/messenger", http.StatusSeeOther)
}

func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	port := utils.GetPort(r)
	sessionName := "session-" + port
	curSession, _ := session.Store.Get(r, sessionName)

	userID, ok := curSession.Values["userID"].(uint)
	if !ok {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		conversationID, err := strconv.ParseUint(r.FormValue("conversation_id"), 10, 32)
		if err != nil {
			http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
			return
		}
		content := r.FormValue("content")

		// Создаем новое сообщение
		message := models.Message{
			SenderID:       userID,
			Content:        content,
			ConversationID: uint(conversationID),                     // Преобразуем строку в uint
			Timestamp:      time.Now().Format("2006-01-02 15:04:05"), // Формат времени
		}

		// Сохраняем сообщение в базе данных
		if err := db.Create(&message).Error; err != nil {
			log.Println("Ошибка при сохранении сообщения:", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		// Перенаправляем обратно на страницу переписки
		http.Redirect(w, r, fmt.Sprintf("/conversation?id=%s", conversationID), http.StatusSeeOther)
		return
	}

	// Если метод не POST, перенаправляем на мессенджер
	http.Redirect(w, r, "/messenger", http.StatusSeeOther)
}
