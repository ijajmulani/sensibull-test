package models

type Plan struct {
	ID       string  `gorm:"primary_key" "type:varchar(50)" json:"id" "UNIQUE"`
	Validity int16   `gorm: json:"validity"`
	Cost     float32 `gorm: json:"cost"`
}

func (u *Plan) TableName() string {
	return "plan"
}
