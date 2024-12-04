package auth

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"html/template"
	"itsm/models"
	"log"
	"net/http"
)

var db *gorm.DB

type user models.User

func SetupRoutes(database *gorm.DB) {
	db = database
	http.HandleFunc("/", authHandler)
	http.HandleFunc("/register", registerHandler)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Проверка пользователя
		var user user
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
				return
			}
			log.Println("Ошибка при выполнении запроса:", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		// Проверка пароля
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
			return
		}

		// Успешная авторизация
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/auth/auth.html"))
	if err := tmpl.Execute(w, map[string]interface{}{"Register": false}); err != nil {
		log.Println("Ошибка при выполнении шаблона:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Проверка на существование пользователя
		var count int64
		if err := db.Model(&user{}).Where("username = ?", username).Count(&count).Error; err != nil {
			log.Println("Ошибка при выполнении запроса:", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		if count > 0 {
			http.Error(w, "Пользователь с таким логином уже существует", http.StatusConflict)
			return
		}

		// Хеширование пароля
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Ошибка при хешировании пароля:", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		// Добавление нового пользователя
		newUser := user{Username: username, Password: string(hashedPassword)}
		if err := db.Create(&newUser).Error; err != nil {
			log.Println("Ошибка при добавлении пользователя:", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		// Успешная регистрация
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/auth/auth.html"))
	if err := tmpl.Execute(w, map[string]interface{}{"Register": true}); err != nil {
		log.Println("Ошибка при выполнении шаблона:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}
