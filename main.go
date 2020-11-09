package main

import (
	"fmt"
	"net/http"
	"os"
	"sensibull-test/routers"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	initDB()
	r := routers.SetupRouter()
	http.ListenAndServe(":"+os.Getenv("HTTPPORT"), r)
}

type User struct {
	// gorm.Model
	ID        uint      `json:"id" gorm:"primary_key"`
	Name      string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at":`
}

var db *gorm.DB

func initDB() {
	var err error
	dsn := "docker:docker@tcp(fullstack-mysql:3306)/sensibull?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println(err)
	}

	// Create the database. This is a one-time step.
	// Comment out if running multiple times - You may see an error otherwise
	db.Exec("CREATE DATABASE IF NOT EXISTS sensibull")
	db.Exec("USE sensibull")

	// Migration to create tables for Order and Item schema
	db.AutoMigrate(&User{})
}
