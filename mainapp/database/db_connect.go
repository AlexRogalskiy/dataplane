package database

package gotimescaledb

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func DBConnect() (*gorm.DB, error) {
	// dialect, connectionURL := postgresConnectionURL(config)
	connectURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("secret.db_user"),
		os.Getenv("secret.db_pwd"),
		os.Getenv("secret.db_host"),
		os.Getenv("secret.db_port"),
		os.Getenv("secret.db_database"),
		os.Getenv("secret.db_ssl"),
	)

	// log.Println(connectionURL)
	var l logger.LogLevel
	dbDebug, _ := strconv.ParseBool(os.Getenv("db_debug"))
	if dbDebug {
		l = logger.Info
		log.Println("DB logging: Info")
	} else {
		l = logger.Silent
	}

	dbConn, err := gorm.Open(postgres.New(postgres.Config{
		DSN: connectURL,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
		Logger: logger.Default.LogMode(l),
	})

	if err != nil {
		return nil, err
	}

	//--- Connection pooling
	sqlDB, _ := dbConn.DB()
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(25)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(25)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	log.Println("db: connected")

	return dbConn, nil
}