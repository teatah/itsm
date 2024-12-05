// main.go
package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"itsm/api/auth"
	"itsm/api/dashboard"
	"itsm/models"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func startServer(port string, db *gorm.DB) {
	r := mux.NewRouter()
	auth.SetupRoutes(r, db)
	dashboard.SetupRoutes(r, db)

	fs := http.FileServer(http.Dir("./templates"))
	r.PathPrefix("/templates/").Handler(http.StripPrefix("/templates/", fs))

	log.Printf("Сервер запущен на :%s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
	}
}

func checkEnvVariables() {
	requiredEnvVars := []string{"DB_USER",
		"DB_HOST",
		"DB_NAME",
		"PORT1",
		"PORT2",
		"USE_2_SERVERS"}

	for _, envVar := range requiredEnvVars {
		value, exists := os.LookupEnv(envVar)
		if !exists || value == "" {
			log.Fatalf("Error: Environment variable %s is required but not set or empty.", envVar)
		}
	}
}

func getEnv(name string) string {
	return os.Getenv(name)
}

func startGoroutines(db *gorm.DB) {
	go startServer(getEnv("PORT1"), db)

	use2servers, err := strconv.ParseBool(getEnv("USE_2_SERVERS"))
	if err != nil {
		log.Fatalf("Error converting .env var USE_2_SERVERS to boolean: %v", err)
	}

	if use2servers {
		go startServer(getEnv("PORT1"), db)
	}

	select {}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	checkEnvVariables()

	dbUser := getEnv("DB_USER")
	dbPass := getEnv("DB_PASS")
	dbHost := getEnv("DB_HOST")
	dbName := getEnv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPass, dbHost, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.User{}, &models.ServiceLine{}, &models.Service{})
	if err != nil {
		log.Fatal(err)
	}

	startGoroutines(db)
}
