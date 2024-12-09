package utils

import (
	"encoding/json"
	"github.com/gorilla/sessions"
	"itsm/session"
	"net/http"
	"strings"
)

func GetCurSession(r *http.Request) (*sessions.Session, error) {
	hostParts := strings.Split(r.Host, ":")
	var port string
	if len(hostParts) > 1 {
		port = hostParts[1]
	} else {
		port = "8080"
	}

	sessionName := "session-" + port
	curSession, err := session.Store.Get(r, sessionName)
	if err != nil {
		return nil, err
	}

	return curSession, nil
}

func GetCurUserID(w http.ResponseWriter, r *http.Request) (uint, error) {
	curSession, err := GetCurSession(r)
	if err != nil {
		http.Error(w, "Ошибка получения сессии: "+err.Error(), http.StatusUnauthorized)
		return 0, err
	}

	userID := curSession.Values["userID"]

	var curUserID uint
	if value, ok := userID.(uint); ok {
		curUserID = value
	} else {
		http.Error(w, "ID текущего пользователя не найден", http.StatusNotFound)
	}

	return curUserID, err
}

func SendJSON(w http.ResponseWriter, dataStruct interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(dataStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func IsClientUser(r *http.Request) (bool, error) {
	curSession, err := GetCurSession(r)
	if err != nil {
		return false, err
	}

	isAdmin := curSession.Values["isAdmin"].(bool)
	isTechOfficer := curSession.Values["isTechOfficer"].(bool)
	isDefaultOfficer := curSession.Values["isDefaultOfficer"].(bool)

	return !isAdmin && !isTechOfficer && !isDefaultOfficer, nil
}
