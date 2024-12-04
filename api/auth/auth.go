package auth

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
)

var db *sql.DB

func SetupRoutes(database *sql.DB) {
	db = database
	http.HandleFunc("/", authHandler)
	http.HandleFunc("/register", registerHandler)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Проверка пользователя
		var storedHash string
		err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedHash)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
				return
			}
			log.Println("Ошибка при выполнении запроса:", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		// Проверка пароля
		err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
		if err != nil {
			http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
			return
		}

		// Успешная авторизация
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/auth.html"))
	tmpl.Execute(w, map[string]interface{}{"Register": false})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Проверка на существование пользователя
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
		if err != nil {
			log.Println("Ошибка при выполнении запроса:", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		if exists {
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
		_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
		if err != nil {
			log.Println("Ошибка при добавлении пользователя:", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		// Успешная регистрация
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/auth.html"))
	tmpl.Execute(w, map[string]interface{}{"Register": true})
}
