package models

import (
	"fmt"
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
	username := "docker"        // os.Getenv("db_user")
	password := "docker"        //os.Getenv("db_pass")
	dbName := "sensibull"       //os.Getenv("db_name")
	dbHost := "fullstack-mysql" // os.Getenv("db_host")
	dbPort := "3306"            // os.Getenv("db_port")

	dsn := username + ":" + password + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Print(err)
	}
	db = conn

	//Automatically create migration as per model
	err = db.AutoMigrate(
		&User{},
		&Plan{},
	)
	var plan Plan
	result := db.First(&plan)
	if result.RowsAffected == 0 {
		var plans = []Plan{
			{ID: "FREE", Validity: -1, Cost: 0.0},
			{ID: "TRIAL", Validity: 7, Cost: 0.0},
			{ID: "LITE_1M", Validity: 30, Cost: 100.0},
			{ID: "PRO_1M", Validity: 30, Cost: 200.0},
			{ID: "LITE_6M", Validity: 180, Cost: 500.0},
			{ID: "PRO_6M", Validity: 180, Cost: 900.0},
		}
		db.Create(&plans)
	}
}

func GetDB() *gorm.DB {
	return db
}
