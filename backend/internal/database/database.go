package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

// Environment indicates which database environment is connected
type Environment string

const (
	EnvLocal Environment = "LOCAL"
	EnvProd  Environment = "PROD"
)

// CurrentEnv stores the current database environment
var CurrentEnv Environment

// Connect establishes a connection pool to the PostgreSQL database
// Checks PROD_DATABASE_URL first, falls back to DATABASE_URL (local)
func Connect() error {
	var databaseURL string

	// Check for production database first
	if prodURL := os.Getenv("PROD_DATABASE_URL"); prodURL != "" {
		databaseURL = prodURL
		CurrentEnv = EnvProd
		log.Println("Connecting to PostgreSQL (PROD)...")
	} else if localURL := os.Getenv("DATABASE_URL"); localURL != "" {
		databaseURL = localURL
		CurrentEnv = EnvLocal
		log.Println("Connecting to PostgreSQL (LOCAL)...")
	} else {
		return fmt.Errorf("no database URL configured (set DATABASE_URL or PROD_DATABASE_URL)")
	}

	var err error
	Pool, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	// Test the connection
	if err := Pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	log.Printf("PostgreSQL connected successfully (%s)", CurrentEnv)
	return nil
}

// Close closes the database connection pool
func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
