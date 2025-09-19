package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/user-reward/internal/config"
	postgress "github.com/user-reward/internal/database/postgres"
	"github.com/user-reward/internal/repository"

	_ "github.com/lib/pq"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	rewards  *postgress.RewardRepository
}

func main() {
	addr := flag.String("addr", ":4000", "Сетевой адрес HTTP")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Println(cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
		cfg.PostgresSSLMode)

	// Подключение к БД
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.SetMaxOpenConns(100)
	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		rewards:  &postgress.RewardRepository{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Запуск сервера на %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
