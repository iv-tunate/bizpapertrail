package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("[FATAL ERROR] PORT env variable not set")
	}

	db_url := os.Getenv("DB_URL")
	if db_url == ""{
		log.Fatal("[FATAL ERROR] Database url env variable not set")
	}

	db_conn, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalf("[FATAL ERROR] Failure to open database connection... ERROR DETAILS: %v", err)
	}
	defer db_conn.Close()

	e := echo.New()
	server := &http.Server{
		Handler: e,
		Addr: ":" + port,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	e.GET("/healthz", checkserverstatus)
	registerRoutes(e)
	//------------------------------------------------------------------------------------
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func(){
		<-shutdown
		log.Println("Shutdown signal received")

		ctx, cancel := context.WithTimeout(context.Background(), 7 * time.Second)
		defer cancel()

		if err = server.Shutdown(ctx); err != nil{
			log.Fatalf("[FATAL ERROR]: Server forcefully shutdown: %v", err)
		}
	}()
	log.Printf("starting up bizpapertrail server on PORT:%v \n\n", port)
	
	if err = server.ListenAndServe(); err != http.ErrServerClosed{
		log.Fatalf("[FATAL ERROR] :Server crashed%v", err)
	}
	log.Println("Gracefully Shutting down bizpapertrail server")
}
