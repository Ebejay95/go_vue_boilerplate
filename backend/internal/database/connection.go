package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "time"

    _ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        dbURL = "postgres://postgres:password@localhost:5432/myapp?sslmode=disable"
    }

    var db *sql.DB
    var err error

    maxRetries := 30
    for i := 0; i < maxRetries; i++ {
        db, err = sql.Open("postgres", dbURL)
        if err != nil {
            log.Printf("Attempt %d: Failed to open database connection: %v", i+1, err)
            time.Sleep(1 * time.Second)
            continue
        }

        err = db.Ping()
        if err != nil {
            log.Printf("Attempt %d: Failed to ping database: %v", i+1, err)
            db.Close()
            time.Sleep(1 * time.Second)
            continue
        }

        log.Printf("Successfully connected to database after %d attempts", i+1)
        return db, nil
    }

    return nil, fmt.Errorf("failed to connect to database after %d attempts: %v", maxRetries, err)
}