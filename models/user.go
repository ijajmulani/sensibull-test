package models

type User struct {
	Model
	Name string `gorm:"unique" "type:varchar(255)" json:"name"` // make unique key
}

func (u *User) TableName() string {
	return "user"
}
