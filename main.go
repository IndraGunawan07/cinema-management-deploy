package main

import (
	"cinema-management/database"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

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

	// 1. Immediate startup logging
	log.Println("🚀 Application booting")
	startTime := time.Now()

	// 2. Initialize Gin in release mode
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// 3. Pre-start health check (for Railway)
	router.GET("/preflight", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// 4. Database connection with retries
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("❌ DATABASE_URL not set")
	}

	var DB *sql.DB
	var err error
	for i := 0; i < 5; i++ {
		DB, err = sql.Open("postgres", dbURL)
		if err == nil {
			err = DB.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("🔁 DB attempt %d: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("💥 DB connection failed:", err)
	}
	defer DB.Close()
	log.Println("✅ Database connected")

	// 5. Run migrations
	log.Println("🔄 Running migrations...")
	database.DBMigrate(DB)
	log.Println("👍 Migrations completed")

	// 6. Real health check
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "healthy",
			"uptime":   time.Since(startTime).String(),
			"database": "connected",
		})
	})

	// 7. Start server in goroutine with keep-alive
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: router,
	}

	// Run server in goroutine
	go func() {
		log.Printf("🌐 Server listening on :%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("💥 Server failed: %v", err)
		}
	}()

	// Block forever
	select {}
}
