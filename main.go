// main.go
package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"itsm/api/auth"
	"itsm/api/dashboard"
	"itsm/models"
	"log"
	"net/http"
)

func main() {
	dsn := "root:@tcp(localhost:3306)/itsm"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(err)
	}

	auth.SetupRoutes(db)
	dashboard.SetupRoutes(db)

	log.Println("Сервер запущен на :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
	}
}
