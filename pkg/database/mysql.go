package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"time"
)

var databasePool *sql.DB
var gormDB *gorm.DB

func SetUpMysqlConnection(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("mysql", os.Getenv("DB_DSN"))
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(25)
	db.SetMaxOpenConns(25)
	db.SetConnMaxLifetime(10 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	err = openDbConnection(ctx, db)

	if err != nil {
		return nil, err
	}

	gormConnection, err := gorm.Open(gormMysql.New(gormMysql.Config{
		Conn: db,
	}), &gorm.Config{})

	gormDB = gormConnection
	databasePool = db

	return db, nil
}

func GetDBPool() *sql.DB {
	return databasePool
}

func GetGormDB() *gorm.DB {
	return gormDB
}

func openDbConnection(ctx context.Context, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}

	return nil
}
