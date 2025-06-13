package main

import (
	"cinema-management/controllers"
	"cinema-management/database"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

var (
	DB  *sql.DB
	err error
)

func main() {
	log.Println("1. Application starting") // First line in main()
	gin.SetMode(gin.ReleaseMode)
	log.Println("2. Gin mode set")

	// err := godotenv.Load("config/.env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	// psqlInfo := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
	// 	os.Getenv("PGHOST"),
	// 	os.Getenv("PGPORT"),
	// 	os.Getenv("PGUSER"),
	// 	os.Getenv("PGPASSWORD"),
	// 	os.Getenv("PGDATABASE"),
	// )

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("DATABASE_URL not set, using manual config")
		dbURL = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
			os.Getenv("PGHOST"),
			os.Getenv("PGPORT"),
			os.Getenv("PGUSER"),
			os.Getenv("PGPASSWORD"),
			os.Getenv("PGDATABASE"),
		)
	}

	DB, err = sql.Open("postgres", dbURL)

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
	router.SetTrustedProxies(nil)
	router.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	router.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})
	router.GET("/cinemas", controllers.GetaLLCinema)
	router.POST("/cinemas", controllers.InsertCinema)
	router.PUT("/cinemas/:id", controllers.UpdateCinema)
	router.DELETE("/cinemas/:id", controllers.DeleteCinema)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Environment variables:")
	for _, e := range os.Environ() {
		log.Println(e)
	}

	// Start server with proper error handling
	log.Println("Starting server at port:", port)
	err = router.Run("0.0.0.0:" + port)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
	log.Println("Server has stopped unexpectedly")
	// fmt.Println("Successfully connected!")
}
