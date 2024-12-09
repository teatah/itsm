package auth

import (
	"errors"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"html/template"
	"itsm/models"
	_ "itsm/session"
	"itsm/utils"
	"log"
	"net/http"
)

var db *gorm.DB
var serverErrorText = "Ошибка сервера. Попробуйте позже"

type user models.User

func SetupRoutes(r *mux.Router, database *gorm.DB) {
	db = database
	r.HandleFunc("/", authHandler)
	r.HandleFunc("/register", registerHandler)
	r.HandleFunc("/logout", logoutHandler)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	errorMessage := ""

	if r.Method == http.MethodPost {
		user, errorMessage := authUser(r)

		// Успешная авторизация
		if len(errorMessage) == 0 {
			curSession, err := utils.GetCurSession(r)
			if err != nil {
				http.Error(w, "Ошибка получения сессии: "+err.Error(), http.StatusUnauthorized)
			}

			curSession.Values["userID"] = user.ID
			curSession.Values["isAdmin"] = user.IsAdmin
			curSession.Values["isTechOfficer"] = user.IsTechOfficer
			curSession.Values["isDefaultOfficer"] = user.IsDefaultOfficer // Сохраняем права доступа
			err = curSession.Save(r, w)
			if err != nil {
				http.Error(w, "Ошибка сохранения сессии", http.StatusUnauthorized)
			}

			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/auth/auth.html"))
	err := tmpl.Execute(w, map[string]interface{}{
		"Register":     false,
		"ErrorMessage": errorMessage})

	if err != nil {
		log.Println("Ошибка при выполнении шаблона:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}

func authUser(r *http.Request) (user, string) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	var user user
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, "Неверный логин или пароль"
		}
		log.Println("Ошибка при выполнении запроса:", err)
		return user, serverErrorText
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, "Неверный логин или пароль"
	}
	return user, ""
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	errorMessage := ""
	if r.Method == http.MethodPost {
		errorMessage = registerUser(r)

		if len(errorMessage) == 0 {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/auth/auth.html"))
	err := tmpl.Execute(w, map[string]interface{}{
		"Register":     true,
		"ErrorMessage": errorMessage})

	if err != nil {
		log.Println("Ошибка при выполнении шаблона:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}

func registerUser(r *http.Request) string {
	username := r.FormValue("username")
	password := r.FormValue("password")

	var count int64
	if err := db.Model(&user{}).Where("username = ?", username).Count(&count).Error; err != nil {
		log.Println("Ошибка при выполнении запроса:", err)
		return serverErrorText
	}

	if count > 0 {
		return "Пользователь с таким логином уже существует"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Ошибка при хешировании пароля:", err)
		return serverErrorText
	}

	newUser := user{Username: username, Password: string(hashedPassword)}
	if err := db.Create(&newUser).Error; err != nil {
		log.Println("Ошибка при добавлении пользователя:", err)
		return serverErrorText
	}
	return ""
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	curSession, err := utils.GetCurSession(r)
	if err != nil {
		http.Error(w, "Ошибка получения сессии: "+err.Error(), http.StatusUnauthorized)
		return
	}

	curSession.Options.MaxAge = -1

	err = curSession.Save(r, w)
	if err != nil {
		http.Error(w, "Ошибка при выходе", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
