package database

import (
	"dataplane/database/models"
	"fmt"
	"log"
	"os"
	"strconv"

	// "gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Migrate() {

	connectURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("secret.db_user"),
		os.Getenv("secret.db_pwd"),
		os.Getenv("secret.db_host"),
		os.Getenv("secret.db_port"),
		os.Getenv("secret.db_database"),
		os.Getenv("secret.db_ssl"),
	)

	var l logger.LogLevel
	dbDebug, _ := strconv.ParseBool(os.Getenv("dbdebug"))
	if dbDebug {
		l = logger.Info
		log.Println("DB logging: Info")
	} else {
		l = logger.Silent
		log.Println("DB logging: Silent")
	}

	dbConn, err := gorm.Open(postgres.Open(connectURL), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
		Logger: logger.Default.LogMode(l),
	})
	if err != nil {
		panic(err)
	}

	err1 := dbConn.AutoMigrate(
		&models.Pipelines{},
		&models.Users{},
		&models.Claims{},
	)
	if err1 != nil {
		panic(err1)
	}

	// ---- load in mcc
	// dbConn.Clauses(clause.OnConflict{
	// 	UpdateAll: true,
	// }).Create(&dataingest.Mcctradeingest)
	// log.Println("mcc loaded")

	// // ---- load country data
	// dbConn.Clauses(clause.OnConflict{
	// 	UpdateAll: true,
	// }).Create(&dataingest.CountryPrices)
	// log.Println("country default pricing loaded")

	// hypertable := "SELECT create_hypertable('logs_billing', 'created_at', if_not_exists => TRUE, chunk_time_interval=> INTERVAL '7 Days');"

	// if hypertable != "" {
	// 	if err := dbConn.Model(&gushbilling.LogsBilling{}).Exec(hypertable).Error; err != nil {
	// 		panic(err)
	// 	}
	// }
	log.Println("Database migrated")
}
