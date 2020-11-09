package models

type User struct {
	Model
	Name string `gorm:"type:varchar(50)" json:"name" validate:"required"`
}

func (u *User) TableName() string {
	return "user"
}
