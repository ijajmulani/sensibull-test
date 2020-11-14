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
	var newStartDate time.Time
	var err error
	db := models.GetDB()

	var user models.User
	userDBRes := db.Where("name = ?", args.UserName).First(&user)
	if userDBRes.RowsAffected == 0 {
		return errors.New("user does not exists")
	}

	var plan models.Plan
	planDBRes := db.Where("name = ?", args.PlanName).First(&plan)
	if planDBRes.RowsAffected == 0 {
		return errors.New("plan does not exists")
	}

	// start should be valid and start date should be greater than current date

	// remaining match date only not time
	if newStartDate, err = time.Parse(layoutISO, args.StartDate); err != nil || newStartDate.Before(time.Now()) {
		return errors.New("start_date is not valid")
	}

	// future start date should not be on overlap

	// check if any plan exists on given date

	// handle free tier case bcz validity is -1
	//select user.id, user.name, subscription.plan_id, subscription.start_date, validity  from user join subscription on user.id = subscription.user_id  join plan on plan.id = subscription.plan_id where
	// '2020-12-13' between start_date and start_date + interval validity day

	// cant upgrade/ degrade plan on current date

	var subscription models.Subscription
	db.Debug().Last(&subscription).Where("user_id = ?", user.ID).Order("start_date")

	if subscription.PlanID != 0 {
		if newStartDate.After(subscription.StartDate) || newStartDate.Equal(subscription.StartDate) {
			if newStartDate.Before(subscription.ValidTill) {
				// Update previous plan's valid_till date
				if newStartDate.Equal(subscription.StartDate) == false || (newStartDate.Equal(subscription.StartDate) && plan.ID != subscription.PlanID) {
					db.Model(&models.Subscription{}).Where("id = ?", subscription.ID).Update("valid_till", newStartDate)
				} else {
					return errors.New("plan is already activated at give date. please choose another plan or provide future date")
				}

			}
		} else {
			return errors.New("start_date is not valid")
		}
	}

	//  hit payment api then according write in db
	newSubscription := models.Subscription{
		PlanID:    plan.ID,
		UserID:    user.ID,
		StartDate: newStartDate,
		ValidTill: newStartDate.AddDate(0, 0, int(plan.Validity)),
	}

	db.Create(&newSubscription)

	return nil
}
