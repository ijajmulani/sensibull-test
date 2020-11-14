package models

type Plan struct {
	Model
	Name     string  `"type:varchar(255)" json:"name"`
	Validity int16   `gorm: json:"validity"` // -1 for infinite
	Cost     float32 `gorm: json:"cost"`     //decimal(10,2) NOT NULL
	// is active for plan update
}

func (u *Plan) TableName() string {
	return "plan"
}
