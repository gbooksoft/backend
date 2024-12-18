package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

var databasePool *sql.DB

func SetUpMysqlConnection(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("mysql", os.Getenv("DB_DSN"))
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(3)
	db.SetMaxOpenConns(3)

	err = openDbConnection(ctx, db)

	if err != nil {
		return nil, err
	}

	databasePool = db

	return db, nil
}

func GetDBPool() *sql.DB {
	return databasePool
}

func openDbConnection(ctx context.Context, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}

	return nil
}
