package storage

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(ctx context.Context, maxAttempts int) (*pgxpool.Pool, error) {
	fmt.Println("Try to connect to db")

	// Read config from env (use docker service name 'postgres' by default)
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres_user"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "12345"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "avitotest_db"
	}

	// sensible defaults
	dbConnectionTimeOut := 5
	if v, err := strconv.Atoi(os.Getenv("DB_CONNECTION_TIMEOUT")); err == nil && v > 0 {
		dbConnectionTimeOut = v
	}

	dbConnectionRetries := maxAttempts
	if v, err := strconv.Atoi(os.Getenv("DB_CONNECTION_RETRYES")); err == nil && v > 0 {
		dbConnectionRetries = v
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	fmt.Println("Dsn =", dsn)

	// try connect with retries
	for attempt := 1; attempt <= dbConnectionRetries; attempt++ {
		fmt.Printf("Try %d/%d\n", attempt, dbConnectionRetries)

		ctx2, cancel := context.WithTimeout(ctx, time.Duration(dbConnectionTimeOut)*time.Second)
		dbpool, err := pgxpool.New(ctx2, dsn)
		if err == nil {
			err = dbpool.Ping(ctx2)
		}
		if err == nil {
			cancel()
			fmt.Println("DB connected!")
			return dbpool, nil
		}

		// cleanup and wait before retry
		if dbpool != nil {
			dbpool.Close()
		}
		cancel()
		fmt.Printf("Attempt %d: Failed to connect: %v\n", attempt, err)
		time.Sleep(time.Second)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts", dbConnectionRetries)
}

func MakeTeam(pool *pgxpool.Pool, ctx context.Context) error {
	// Начинаем транзакцию
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // Всегда откатываем при ошибке

	// Выполняем операции в транзакции
	_, err = tx.Exec(ctx, "select team_name from teams", 100, 1)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2", 100, 2)
	if err != nil {
		return err
	}

	// Коммитим транзакцию
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
