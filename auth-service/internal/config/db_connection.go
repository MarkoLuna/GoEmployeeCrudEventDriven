package config

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func InitDatabase() {
	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbPort := utils.GetEnv("DB_PORT", "5432")
	dbName := utils.GetEnv("DB_NAME", "auth_db")
	dbUser := utils.GetEnv("DB_USER", "auth_user")
	dbPassword := utils.GetEnv("DB_PASSWORD", "authpw")
	dbDriverName := utils.GetEnv("DB_DRIVER_NAME", "postgres")

	connectDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword)

	var conn *sql.DB
	var err error
	for i := 0; i < 30; i++ {
		conn, err = sql.Open(dbDriverName, connectDSN)
		if err == nil {
			err = conn.Ping()
		}
		if err == nil {
			break
		}
		log.Printf("Waiting for database... (%d/30): %v", i+1, err)
		time.Sleep(time.Second)
	}
	if err != nil {
		log.Fatalf("Could not connect to database after 30 retries: %v", err)
	}
	defer conn.Close()

	var exists bool
	err = conn.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if err != nil {
		log.Fatalf("Could not check database existence: %v", err)
	}

	if !exists {
		_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			log.Fatalf("Could not create database %s: %v", dbName, err)
		}
		log.Printf("Created database: %s", dbName)
	}

	targetDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	targetDB, err := sql.Open(dbDriverName, targetDSN)
	if err != nil {
		log.Fatalf("error connecting to target DB: %v", err)
	}
	defer targetDB.Close()

	if err := targetDB.Ping(); err != nil {
		log.Fatalf("error pinging target DB: %v", err)
	}

	runMigrations(targetDB, dbDriverName)
}

func runMigrations(db *sql.DB, driverName string) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not create migration driver: %v", err)
	}

	src, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		log.Fatalf("Could not create migration source: %v", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, driverName, driver)
	if err != nil {
		log.Fatalf("Could not create migration instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Could not run migrations: %v", err)
	}

	log.Println("Database migrations completed successfully")
}

func GetDB() *sql.DB {
	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbPort := utils.GetEnv("DB_PORT", "5432")
	dbName := utils.GetEnv("DB_NAME", "auth_db")
	dbUser := utils.GetEnv("DB_USER", "auth_user")
	dbPassword := utils.GetEnv("DB_PASSWORD", "authpw")
	dbDriverName := utils.GetEnv("DB_DRIVER_NAME", "postgres")

	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open(dbDriverName, dataSourceName)
	if err != nil {
		log.Fatal("error connection to DB: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("error pinging DB: ", err)
	}

	return db
}
