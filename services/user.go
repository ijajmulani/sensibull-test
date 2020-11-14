package services

import (
	"errors"
	"sensibull-test/models"
)

type UserService struct {
}
type UserResponse struct {
	Name      string `json:"user_name"`
	CreatedAt string `json:"created_at"`
}

// Get will get user details from db
func (us *UserService) Get(userName string) (UserResponse, error) {
	var user models.User
	db := models.GetDB()
	result := db.Where("name = ?", userName).First(&user)
	if result.RowsAffected == 0 {
		return UserResponse{}, errors.New("user does not exists")
	}
	layoutISO := "2006-01-02"
	res := UserResponse{
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format(layoutISO),
	}
	return res, nil
}

// Add will add user details in db
func (us *UserService) Add(name string) error {
	db := models.GetDB()
	var user models.User
	result := db.Where("name = ?", name).First(&user)
	if result.RowsAffected > 0 {
		return errors.New("user_name already exists")
	}
	user.Name = name
	result = db.Create(&user)
	return result.Error
}
