package models

import "time"

type Subscription struct {
	Model
	Plan      Plan
	User      User
	StartDate time.Time `json:"start_date"`
	ValidTill time.Time `json:"valid_till`
	PlanID    uint
	UserID    uint
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Subscription) TableName() string {
	return "subscription"
}
