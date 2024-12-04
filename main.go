// main.go
package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"itsm/api" // Замените на имя вашего модуля
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/itsm")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создание таблицы пользователей
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        username VARCHAR(255) NOT NULL,
        password VARCHAR(255) NOT NULL,
        is_admin BOOLEAN DEFAULT FALSE
    )`)
	if err != nil {
		log.Fatal(err)
	}

	api.SetupRoutes(db)

	log.Println("Сервер запущен на :8080")
	http.ListenAndServe(":8080", nil)
}
