package main

import (
	"database/sql"
	"itsm/api/auth"
	"itsm/api/dashboard"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/dbname")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создание таблицы пользователей с новой колонкой is_admin
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        username VARCHAR(255) NOT NULL,
        password VARCHAR(255) NOT NULL,
        is_admin BOOLEAN DEFAULT FALSE
    )`)
	if err != nil {
		log.Fatal(err)
	}

	auth.SetupRoutes(db)
	dashboard.SetupRoutes(db)

	log.Println("Сервер клиентов запущен на :8081")
	http.ListenAndServe(":8081", nil)
}
