package main

import (
	"flag"
	"os"
	"time"
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
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development|production|statging)")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("DSN"), "PostgreSQL DSN")

	flag.Parse()

	app := &application{
		config: cfg,
	}

	app.serve()
}
