package services

import (
	"errors"
	"sensibull-test/models"
	"sensibull-test/structures/subscriptions"
	"time"
)

type SubscriptionService struct {
}

func (ss *SubscriptionService) Post(args subscriptions.PostArgs) error {
	const layoutISO = "2006-01-02"
	var startDate time.Time
	var err error
	db := models.GetDB()
	var user models.User
	userDBRes := db.Where("name = ?", args.UserName).First(&user)
	if userDBRes.RowsAffected == 0 {
		return errors.New("user does not exists")
	}

	var plan models.Plan
	planDBRes := db.Where("id = ?", args.PlanID).First(&plan)
	if planDBRes.RowsAffected == 0 {
		return errors.New("plan does not exists")
	}

	// start should be valid and start date should be greater than current date
	if startDate, err = time.Parse(layoutISO, args.StartDate); err != nil || startDate.Before(time.Now()) {
		return errors.New("start_date is not valid")
	}

	// future start date should not be on overlap

	// check if any plan exists on given date

	//select user.id, user.name, subscription.plan_id, subscription.start_date, validity  from user join subscription on user.id = subscription.user_id  join plan on plan.id = subscription.plan_id where
	// '2020-12-13' between start_date and start_date + interval validity day

	subscription := models.Subscription{
		StartDate: startDate,
		PlanID:    args.PlanID,
		UserID:    user.ID,
	}

	db.Create(&subscription)
	return nil
}
