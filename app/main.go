package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"com.derso/testecargak6/amigos"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	gormDB, sqlDB := connectDB()
	defer sqlDB.Close()
	server := createServer(gormDB)
	setupGracefulShutdown(server)
}

func connectDB() (*gorm.DB, *sql.DB) {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port

	logLevel := logger.Silent
	envLevel := os.Getenv("DB_LOG_LEVEL")

	switch envLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logLevel,
			Colorful:      true,
		},
	)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	sqlDB, err := gormDB.DB()

	if err != nil {
		log.Fatal("failed to get underlying sql.DB", err)
	}

	// OTIMIZADO: pool de conexões
	// Ponto ideal para os recursos finitos do container
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)
	monitorarPool(sqlDB)

	return gormDB, sqlDB
}

func createServer(DB *gorm.DB) *http.Server {
	router := gin.Default()

	amigos.Inicializar(DB, router)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Erro ao criar servidor")
			panic(err)
		}
	}()

	fmt.Println("Servidor criado.")

	return srv
}

func setupGracefulShutdown(server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, os.Interrupt) // os.Interrupt: Ctrl+C
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Println("Parando...")

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Erro ao encerrar servidor:", err)
	}

	<-ctx.Done()
}

func monitorarPool(sqlDB *sql.DB) {
	go func() {
		for {
			stats := sqlDB.Stats()

			log.Printf("📊 [Pool Stats] Abertas: %d | Em Uso: %d | Ociosas: %d | Esperando: %d",
				stats.OpenConnections,
				stats.InUse,
				stats.Idle,
				stats.WaitCount, // Total de requisições que tiveram que esperar por uma conexão (O MAIS IMPORTANTE!)
			)

			time.Sleep(5 * time.Second)
		}
	}()
}
