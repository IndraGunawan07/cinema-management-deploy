package main

import (
	"cinema-management/database"
	"database/sql"
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
	// log.Println("1. Application starting") // First line in main()
	// gin.SetMode(gin.ReleaseMode)
	// log.Println("2. Gin mode set")

	// // err := godotenv.Load("config/.env")
	// // if err != nil {
	// // 	log.Fatal("Error loading .env file")
	// // }

	// // psqlInfo := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
	// // 	os.Getenv("PGHOST"),
	// // 	os.Getenv("PGPORT"),
	// // 	os.Getenv("PGUSER"),
	// // 	os.Getenv("PGPASSWORD"),
	// // 	os.Getenv("PGDATABASE"),
	// // )

	// dbURL := os.Getenv("DATABASE_URL")
	// if dbURL == "" {
	// 	log.Println("DATABASE_URL not set, using manual config")
	// 	dbURL = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
	// 		os.Getenv("PGHOST"),
	// 		os.Getenv("PGPORT"),
	// 		os.Getenv("PGUSER"),
	// 		os.Getenv("PGPASSWORD"),
	// 		os.Getenv("PGDATABASE"),
	// 	)
	// }

	// DB, err = sql.Open("postgres", dbURL)

	// if err != nil {
	// 	log.Fatalf("Error opening database: %v\n", err)
	// }

	// defer DB.Close()
	// err = DB.Ping()
	// if err != nil {
	// 	panic(err)
	// }

	// database.DBMigrate(DB)

	// router := gin.Default()
	// router.SetTrustedProxies([]string{"127.0.0.1"})
	// router.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	// router.GET("/cinemas", controllers.GetaLLCinema)
	// router.POST("/cinemas", controllers.InsertCinema)
	// router.PUT("/cinemas/:id", controllers.UpdateCinema)
	// router.DELETE("/cinemas/:id", controllers.DeleteCinema)

	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "8080"
	// }

	// // Start server with proper error handling
	// if err := router.Run("0.0.0.0:" + port); err != nil {
	// 	panic("Server failed to start: " + err.Error())
	// }
	// // fmt.Println("Successfully connected!")

	log.Println("1. Application starting")
	gin.SetMode(gin.ReleaseMode)
	log.Println("2. Gin mode set")

	// Database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	DB, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("DB open failed:", err)
	}
	defer DB.Close()

	// Retry connection
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		err = DB.Ping()
		if err == nil {
			break
		}
		log.Printf("DB ping attempt %d failed: %v", i+1, err)
	}
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	log.Println("3. Database connected")

	// Run migrations
	log.Println("4. Running migrations...")
	database.DBMigrate(DB)
	log.Println("5. Migrations completed")

	// Router setup
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server (with keep-alive)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("6. Starting server on port", port)
	done := make(chan bool)
	go func() {
		if err := router.Run("0.0.0.0:" + port); err != nil {
			log.Fatal("Server error:", err)
		}
	}()
	<-done // Block indefinitely
}
