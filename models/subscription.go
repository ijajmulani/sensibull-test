package models

import "time"

type Subscription struct {
	Model
	StartDate time.Time `json:"start_date"`
	PlanID    uint
	Plan      Plan
	UserID    uint
	User      User
}

func (s *Subscription) TableName() string {
	return "subscription"
}
