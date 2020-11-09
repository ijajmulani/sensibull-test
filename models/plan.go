package models

type Plan struct {
	Model
	Name     string `gorm:"type:varchar(50)" json:"name" validate:"required, UNIQUE"`
	Validity int16  `gorm: json:"validity" validate:"required"`
	Cost     uint16 `gorm: json:"cost" validate:"required"`
}

func (u *Plan) TableName() string {
	return "plan"
}
