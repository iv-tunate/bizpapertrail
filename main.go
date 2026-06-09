package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iv-tunate/bizpapertrail/database"
	"github.com/iv-tunate/bizpapertrail/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
)

func main() {
	godotenv.Load(".env")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	env := os.Getenv("APP_ENV");
	slog.SetDefault(logger)
	

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("[FATAL ERROR] PORT env variable not set")
	}

	db_url := os.Getenv("DB_URL")
	if db_url == ""{
		log.Fatal("[FATAL ERROR] Database url env variable not set")
	}

	conn_pool, err := pgxpool.New(context.Background(), db_url)
	

	if err != nil {
		log.Fatalf("[FATAL ERROR] Failure to open database connection... ERROR DETAILS: %v", err)
	}
	defer conn_pool.Close()

	db_queries := database.New(conn_pool)

	h := handlers.NEW(db_queries, conn_pool, logger)

	e := echo.New()
	server := &http.Server{
		Handler: e,
		Addr: ":" + port,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	e.GET("/healthz", checkserverstatus)
	registerRoutes(e, h)

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
	slog.Info("bizpapertrail running in [environment] mode", env)
	if err = server.ListenAndServe(); err != http.ErrServerClosed{
		log.Fatalf("[FATAL ERROR] :Server crashed%v", err)
	}
	log.Println("Gracefully Shutting down bizpapertrail server")
}
