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
	gin.SetMode(gin.ReleaseMode)

	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	psqlInfo := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGDATABASE"),
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
	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.GET("/cinemas", controllers.GetaLLCinema)
	router.POST("/cinemas", controllers.InsertCinema)
	router.PUT("/cinemas/:id", controllers.UpdateCinema)
	router.DELETE("/cinemas/:id", controllers.DeleteCinema)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server with proper error handling
	if err := router.Run("0.0.0.0:" + port); err != nil {
		panic("Server failed to start: " + err.Error())
	}
	// fmt.Println("Successfully connected!")
}
