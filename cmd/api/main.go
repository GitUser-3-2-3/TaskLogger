package main

import (
	"TaskLogger/internal/data"
	"context"
	"database/sql"
	"flag"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/lib/pq"
)

/*
	todo 0-> refine the entities
	todo 1-> migrate database to postgres
	todo 2-> change all the queries to accommodate to postgres
	todo 3-> instead of pq driver use pgx driver for database
	todo 4-> write authentication
*/

const version = "1.0.0"

type config struct {
	env  string
	port int
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleTime  time.Duration
		maxIdleConns int
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type backend struct {
	logger zerolog.Logger
	config config
	models data.Models
}

func main() {
	var cfg config
	if err := godotenv.Load(".envrc"); err != nil {
		log.Error().Err(err).Msg("Error loading .envrc file")
	}
	runFlags(&cfg)

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	db, err := connectDB(cfg)
	if err != nil {
		log.Fatal().Msg("Error while connecting to DB:: " + err.Error())
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Err(err).Msg("Error while closing DB")
		}
	}(db)
	logger.Info().Msg("Connection with database established")

	bknd := &backend{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}
	err = bknd.serve()
	if err != nil {
		logger.Fatal().Msg("Error while serving http connection")
	}
}

func runFlags(cfg *config) {
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev | prod)")
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("TASK_LOGGER_DSN"), "MySQL DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "Max Open DB Connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "Max Idle DB Connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "Max Idle DB Time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Limiter max requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 5, "Limiter max burst requests")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiting")

	flag.Parse()
}

func connectDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	} else {
		return db, nil
	}
}
