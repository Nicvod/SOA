package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	DBConn        DBConnConfig
	ServicePort   string
	PublicKeyFile string
}

type DBConnConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func NewDB(cfg *Config) (*sqlx.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBConn.DBHost, cfg.DBConn.DBPort, cfg.DBConn.DBUser, cfg.DBConn.DBPassword, cfg.DBConn.DBName,
	)

	var db *sqlx.DB
	var err error

	maxAttempts := 10
	retryInterval := 15 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		log.Println(connStr)
		db, err = sqlx.Connect("postgres", connStr)
		if err == nil {
			log.Println("Successfully connected to the database!")
			return db, nil
		}

		log.Printf("Attempt %d/%d: Failed to connect to database: %v\n", attempt, maxAttempts, err)
		if attempt < maxAttempts {
			log.Printf("Retrying in %v...\n", retryInterval)
			time.Sleep(retryInterval)
		}
	}
	return nil, fmt.Errorf("failed to connect to database after %d attempts: %v", maxAttempts, err)
}

func NewConfig() (*Config, error) {
	var publicFile, dbNameEnv, dbUserEnv, dbPasswordEnv, dbName, dbUser, dbPassword string
	flag.StringVar(&publicFile, "public_key", "", "path to JWT public key `file`")
	flag.StringVar(&dbNameEnv, "db_name_env", "", "database name env")
	flag.StringVar(&dbUserEnv, "db_user_env", "", "database user env")
	flag.StringVar(&dbPasswordEnv, "db_password_env", "", "database password env")
	dbPort := flag.Int("db_port", 5432, "database port")
	servicePort := flag.Int("service_port", 50051, "service port")
	flag.Parse()
	if publicFile == "" {
		return nil, fmt.Errorf("no private key file provided")
	}
	if dbNameEnv == "" {
		return nil, fmt.Errorf("no database name env provided")
	}
	if dbUserEnv == "" {
		return nil, fmt.Errorf("no database user env provided")
	}
	if dbPasswordEnv == "" {
		return nil, fmt.Errorf("no database password env provided")
	}
	dbName = os.Getenv(dbNameEnv)
	dbUser = os.Getenv(dbUserEnv)
	dbPassword = os.Getenv(dbPasswordEnv)
	if dbName == "" || dbPassword == "" || dbUser == "" {
		return nil, fmt.Errorf("not all database info provided")
	}
	return &Config{
		DBConn: DBConnConfig{
			DBHost:     "post_db",
			DBPort:     fmt.Sprint(*dbPort),
			DBUser:     dbUser,
			DBPassword: dbPassword,
			DBName:     dbName,
		},
		ServicePort:   fmt.Sprint(*servicePort),
		PublicKeyFile: publicFile,
	}, nil
}
