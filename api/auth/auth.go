package auth

import (
	"errors"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"html/template"
	"itsm/models"
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
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	errorMessage := ""
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Проверка пользователя
		var user user
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				errorMessage = "Неверный логин или пароль"
			}
			log.Println("Ошибка при выполнении запроса:", err)
		}

		// Проверка пароля
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			errorMessage = "Неверный логин или пароль"
		}

		// Успешная авторизация
		if len(errorMessage) == 0 {
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

func registerHandler(w http.ResponseWriter, r *http.Request) {
	errorMessage := ""
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Проверка на существование пользователя
		var count int64
		if err := db.Model(&user{}).Where("username = ?", username).Count(&count).Error; err != nil {
			log.Println("Ошибка при выполнении запроса:", err)
			errorMessage = serverErrorText
		}

		if count > 0 {
			errorMessage = "Пользователь с таким логином уже существует"
		}

		// Хеширование пароля
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Ошибка при хешировании пароля:", err)
			errorMessage = serverErrorText
		}

		// Добавление нового пользователя
		newUser := user{Username: username, Password: string(hashedPassword)}
		if err := db.Create(&newUser).Error; err != nil {
			log.Println("Ошибка при добавлении пользователя:", err)
			errorMessage = serverErrorText
		}

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
