package messenger

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"itsm/models"
	"itsm/utils"
	"net/http"
	"time"
)

var db *gorm.DB

type ExtendedDialog struct {
	models.Dialog
	Username string `json:"username"`
	Comp     string `json:"comp"`
	CompID   uint   `json:"comp_id"`
}

type ExtendedMessage struct {
	models.Message
	SenderName string `json:"sender_name"`
}

func SetupRoutes(r *mux.Router, database *gorm.DB) {
	db = database
	r.HandleFunc("/users/get", getUsersHandler)
	r.HandleFunc("/dialogs/get", getDialogsHandler)
	r.HandleFunc("/messages/get/{dialogId:[0-9]+}", getMessagesHandler)
	r.HandleFunc("/messages/send", sendMessageHandler).Methods("POST")
	r.HandleFunc("/dialogs/create", createDialogHandler).Methods("POST")
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	userID, err := utils.GetCurUserID(w, r)
	if err != nil {
		return
	}

	if err := db.Table("users").Select("users.*").
		Joins("LEFT JOIN dialogs ON (users.id = dialogs.user1_id AND dialogs.user2_id = ?)"+
			"OR (users.id = dialogs.user2_id AND dialogs.user1_id = ?)", userID, userID).
		Where("users.id != ?", userID).
		Where("users.is_admin OR users.is_tech_officer OR users.is_default_officer").
		Where("dialogs.id IS NULL").
		Find(&users).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, users)
}

func getDialogsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetCurUserID(w, r)
	if err != nil {
		return
	}

	var dialogs []ExtendedDialog

	query := getDialogsQueryText()

	if err := db.Raw(query, userID, userID).Scan(&dialogs).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, dialogs)
}

func getDialogsQueryText() string {
	return `
SELECT dialogs.*, user1.username as username, user2.username as comp, user2.id as comp_id
FROM dialogs
JOIN users as user1 ON user1.id = dialogs.user1_id
JOIN users as user2 ON user2.id = dialogs.user2_id
WHERE user1.id = ?

UNION

SELECT dialogs.*, user2.username as username, user1.username as comp, user1.id as comp_id
FROM dialogs
JOIN users as user2 ON user2.id = dialogs.user2_id
JOIN users as user1 ON user1.id = dialogs.user1_id
WHERE user2.id = ?
`
}

func getMessagesHandler(w http.ResponseWriter, r *http.Request) {
	lastTimestampStr := r.URL.Query().Get("lastTimestamp")
	var lastTimestamp time.Time
	if lastTimestampStr != "" {
		// Парсим временную метку из строки в формате ISO 8601
		lastTimestamp, _ = time.Parse(time.RFC3339, lastTimestampStr)
	}

	vars := mux.Vars(r)
	dialogId := vars["dialogId"]

	var messages []ExtendedMessage

	query := getMessagesQueryText()

	if err := db.Raw(query, dialogId, lastTimestamp).Scan(&messages).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, messages)
}

func getMessagesQueryText() string {
	return `
SELECT messages.*, sender.username as sender_name
FROM messages
JOIN users as sender ON sender.id = messages.sender_id
WHERE messages.dialog_id = ? and messages.timestamp > ?
ORDER BY timestamp ASC
`
}

func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	var message models.Message

	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := utils.GetCurUserID(w, r)
	if err != nil {
		return
	}
	message.SenderID = userID

	if err := db.Create(&message).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func createDialogHandler(w http.ResponseWriter, r *http.Request) {
	var dialog models.Dialog

	if err := json.NewDecoder(r.Body).Decode(&dialog); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.Create(&dialog).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, dialog)
}
