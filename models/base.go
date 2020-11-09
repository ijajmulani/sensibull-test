package models

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"
)

var db *gorm.DB
var validate *validator.Validate

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

	conn, err := gorm.Open("mysql", username+":"+password+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset=utf8&parseTime=True&loc=Asia%2FKolkata")

	if err != nil {
		fmt.Print(err)
	}
	db = conn

	//Printing query
	db.LogMode(true)

	//Automatically create migration as per model
	db.Debug().AutoMigrate(
		&User{},
	)
}

func GetDB() *gorm.DB {
	return db
}
