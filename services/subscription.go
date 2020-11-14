package services

import (
	"errors"
	"sensibull-test/helper"
	"sensibull-test/models"
	"sensibull-test/structures/subscriptions"
	"time"
)

type SubscriptionService struct{}

type PaymentResp struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
}

type SubscriptionListResponse struct {
	PlanName  string `json:"plan_id"`
	StartDate string `json:"start_date"`
	ValidTill string `json:"valid_till"`
}

type SubscriptionGetResponse struct {
	DaysLeft int    `json:"days_left"`
	PlanName string `json:"plan_id"`
}

func (ss *SubscriptionService) GetByUserName(userName string) ([]SubscriptionListResponse, error) {
	db := models.GetDB()
	var subscription models.Subscription
	var response []SubscriptionListResponse
	layoutISO := "2006-01-02"
	rows, err := db.
		Table(subscription.TableName()).
		Select("subscription.start_date, subscription.valid_till, plan.name").
		Joins("join plan on subscription.plan_id = plan.id").
		Joins("join user on subscription.user_id = user.id").
		Where("user.name = ?", userName).
		Order("start_date asc").
		Rows()

	if err == nil {
		var (
			planName  string
			startDate time.Time
			validTill time.Time
		)
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&startDate, &validTill, &planName); err != nil {
				return response, err
			}

			response = append(response, SubscriptionListResponse{
				StartDate: startDate.Format(layoutISO),
				ValidTill: validTill.Format(layoutISO),
				PlanName:  planName,
			})
		}
	}
	if len(response) == 0 {
		response = make([]SubscriptionListResponse, 0)
	}

	return response, err
}

func (ss *SubscriptionService) GetByUserNameAndDate(userName string, date string) (*SubscriptionGetResponse, error) {
	const layoutISO = "2006-01-02"

	inputDate, err := time.Parse(layoutISO, date)
	if err != nil {
		return nil, err
	}

	db := models.GetDB()
	var subscription models.Subscription
	type Result struct {
		ValidTill time.Time
		Name      string
	}
	var result Result
	db.Debug().
		Table(subscription.TableName()).
		Select("subscription.valid_till, plan.name").
		Joins("join user on subscription.user_id = user.id").
		Joins("join plan on subscription.plan_id = plan.id").
		Where("user.name = ? and subscription.start_date <= ? and subscription.valid_till > ?", userName, date, date).
		Scan(&result)

	var response = new(SubscriptionGetResponse)
	if result.Name != "" {
		response.PlanName = result.Name
		response.DaysLeft = int(result.ValidTill.Sub(inputDate).Hours() / 24)
	}
	return response, nil
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

	var newPlan models.Plan
	planDBRes := db.Where("name = ?", args.PlanName).First(&newPlan)
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
	var amountToProcess = -newPlan.Cost

	if subscription.PlanID != 0 {
		var oldPlanUsesDays float32
		if newStartDate.After(subscription.StartDate) || newStartDate.Equal(subscription.StartDate) {
			if newStartDate.Before(subscription.ValidTill) {
				// Update previous plan's valid_till date
				if newPlan.ID != subscription.PlanID || (newStartDate.Equal(subscription.StartDate) && newPlan.ID != subscription.PlanID) {
					db.Model(&models.Subscription{}).Where("id = ?", subscription.ID).Update("valid_till", newStartDate)
					oldPlanUsesDays = float32(newStartDate.Sub(subscription.StartDate).Hours() / 24)
				} else {
					return errors.New("plan is already activated at give date. please choose another plan or provide future date")
				}

			}
		} else {
			return errors.New("start_date is not valid")
		}
		if newPlan.ID != subscription.PlanID {
			var oldPlanCharges float32
			var oldPlan models.Plan
			oldPlanDetails := db.Where("id = ?", subscription.PlanID).First(&oldPlan)
			if oldPlanDetails.RowsAffected == 0 {
				return errors.New("oops, error occurred")
			}

			if oldPlanUsesDays == 0 {
				oldPlanCharges = 0
			} else {
				oldPlanCharges = (oldPlan.Cost / float32(oldPlan.Validity)) * oldPlanUsesDays
			}
			amountToProcess = (oldPlan.Cost - oldPlanCharges)
		}

		amountToProcess = -(newPlan.Cost - amountToProcess)
	}

	paymentResp := new(PaymentResp)
	if err = helper.ProcessPayment(args.UserName, amountToProcess, paymentResp); err == nil {

		newSubscription := models.Subscription{
			PlanID:    newPlan.ID,
			UserID:    user.ID,
			StartDate: newStartDate,
			ValidTill: newStartDate.AddDate(0, 0, int(newPlan.Validity)),
		}

		db.Create(&newSubscription)
	}

	return err
}
