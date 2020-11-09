package services

import (
	"sensibull-test/models"
)

type UserService struct{}

func (us *UserService) List() interface{} {
	var user models.User
	db := models.GetDB()
	db.First(&user)
	return user
}

func (us *UserService) Add() (interface{}, error) {
	var err error
	db := models.GetDB()

	user := models.User{Name: "Jinzhu"}
	result := db.Create(&user)
	if result.RowsAffected == 1 {
		err = result.Error
	}
	return user, err
}
