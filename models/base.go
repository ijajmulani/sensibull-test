package models

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

type Model struct {
	ID        uint      `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt time.Time `gorm:"not null" json:"created_at" sql:"DEFAULT:CURRENT_TIMESTAMP"`
}

func init() {

	username := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")
	dbHost := os.Getenv("MYSQL_ROOT_HOST")
	dbPort := os.Getenv("MYSQL_PORT")

	for {
		dsn := username + ":" + password + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
		conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("Could not connect to database. Retrying in 2 seconds")
			time.Sleep(2 * time.Second)
		} else {
			db = conn
			break
		}
	}

	//Automatically create migration as per model
	db.AutoMigrate(
		&User{},
		&Plan{},
		&Subscription{},
	)
	var plan Plan
	result := db.First(&plan)
	if result.RowsAffected == 0 {
		var plans = []Plan{
			{Name: "FREE", Validity: -1, Cost: 0.0},
			{Name: "TRIAL", Validity: 7, Cost: 0.0},
			{Name: "LITE_1M", Validity: 30, Cost: 100.0},
			{Name: "PRO_1M", Validity: 30, Cost: 200.0},
			{Name: "LITE_6M", Validity: 180, Cost: 500.0},
			{Name: "PRO_6M", Validity: 180, Cost: 900.0},
		}
		db.Create(&plans)
	}
}

func GetDB() *gorm.DB {
	return db
}
