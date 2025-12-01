package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"
	"time"

	"github.com/kharljhon14/zentrix/internal/data"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
}

type application struct {
	config config
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development|production|statging)")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("DSN"), "PostgreSQL DSN")

	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		log.Printf("Failed to connect to DB %v", err)
		os.Exit(1)
	}

	app := &application{
		config: cfg,
		models: data.NewModels(db),
	}

	app.serve()
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetConnMaxLifetime(cfg.db.maxIdleTime)
	db.SetMaxOpenConns(cfg.db.maxOpenConns)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
