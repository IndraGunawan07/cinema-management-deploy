package main

import (
	"cinema-management/database"
	"database/sql"
	"fmt"
	"log"
	"os"

	"cinema-management/controllers"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

var (
	DB  *sql.DB
	err error
)

func main() {
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	psqlInfo := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	DB, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}

	defer DB.Close()
	err = DB.Ping()
	if err != nil {
		panic(err)
	}

	database.DBMigrate(DB)

	router := gin.Default()
	router.GET("/cinemas", controllers.GetaLLCinema)
	router.POST("/cinemas", controllers.InsertCinema)
	router.PUT("/cinemas/:id", controllers.UpdateCinema)
	router.DELETE("/cinemas/:id", controllers.DeleteCinema)

	router.Run(":8080")
	// fmt.Println("Successfully connected!")
}
