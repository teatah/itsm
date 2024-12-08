package utils

import (
	"encoding/json"
	"itsm/session"
	"net/http"
	"strings"
)

func GetPort(r *http.Request) string {
	hostParts := strings.Split(r.Host, ":")
	var port string
	if len(hostParts) > 1 {
		port = hostParts[1]
	} else {
		port = "8080"
	}
	return port
}

func GetCurUserID(w http.ResponseWriter, r *http.Request) (uint, error) {
	port := GetPort(r)
	sessionName := "session-" + port
	curSession, _ := session.Store.Get(r, sessionName)

	userID := curSession.Values["userID"]

	var err error = nil
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
