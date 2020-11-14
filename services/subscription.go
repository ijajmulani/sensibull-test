package services

import (
	"errors"
	"sensibull-test/constants"
	"sensibull-test/helper"
	"sensibull-test/models"

	"sensibull-test/structures/subscriptions"
	"time"

	"gorm.io/gorm"
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
	PlanName string `json:"plan_id"`
	DaysLeft int    `json:"days_left"`
}

type SubscriptionPostResponse struct {
	Status string  `json:"status"`
	Amount float32 `json:"amount"`
}

// GetByUserName return []SubscriptionListResponse of given user name
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

// GetByUserNameAndDate return *SubscriptionGetResponse of given user name at perticular date
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

// Post will add subscription in db
func (ss *SubscriptionService) Post(args subscriptions.PostArgs) (SubscriptionPostResponse, error) {
	const layoutISO = "2006-01-02"
	var newStartDate time.Time
	var err error
	var res SubscriptionPostResponse
	db := models.GetDB()
	res.Status = constants.FAILURE

	var user models.User
	userDBRes := db.Where("name = ?", args.UserName).First(&user)
	if userDBRes.RowsAffected == 0 {
		return res, errors.New(constants.UserNotFound)
	}

	var newPlan models.Plan
	planDBRes := db.Where("name = ?", args.PlanName).First(&newPlan)
	if planDBRes.RowsAffected == 0 {
		return res, errors.New(constants.PlanNotExist)
	}

	// start should be valid and start date should be greater than current date
	// remaining match date only not time
	if newStartDate, err = time.Parse(layoutISO, args.StartDate); err != nil || newStartDate.Before(time.Now()) {
		return res, errors.New(constants.StartDateNotValid)
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
	err = db.Transaction(func(tx *gorm.DB) error {
		if subscription.PlanID != 0 {
			var oldPlanUsesDays float32
			if newStartDate.After(subscription.StartDate) || newStartDate.Equal(subscription.StartDate) {
				if newStartDate.Before(subscription.ValidTill) {
					// Update previous plan's valid_till date
					if newPlan.ID != subscription.PlanID || (newStartDate.Equal(subscription.StartDate) && newPlan.ID != subscription.PlanID) {
						tx.Model(&models.Subscription{}).Where("id = ?", subscription.ID).Update("valid_till", newStartDate)
						oldPlanUsesDays = float32(newStartDate.Sub(subscription.StartDate).Hours() / 24)
					} else {
						return errors.New(constants.PlanAlreadyActivated)
					}

					if newPlan.ID != subscription.PlanID {
						var oldPlanCharges float32
						var oldPlan models.Plan
						oldPlanDetails := tx.Where("id = ?", subscription.PlanID).First(&oldPlan)
						if oldPlanDetails.RowsAffected == 0 {
							return errors.New(constants.ErrorOccurred)
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
			} else {
				return errors.New(constants.StartDateNotValid)
			}
		}

		paymentResp := new(PaymentResp)

		// if amountToProcess == 0 then no need to hit payment api
		if amountToProcess != 0 {
			err = helper.ProcessPayment(args.UserName, amountToProcess, paymentResp)
		} else {
			amountToProcess = 0
		}
		if err == nil {
			res.Amount = amountToProcess
			res.Status = constants.SUCCESS
			var validTill time.Time
			if newPlan.Validity == -1 {
				validTill = newStartDate.AddDate(100, 0, 0)
			} else {
				validTill = newStartDate.AddDate(0, 0, int(newPlan.Validity))
			}

			newSubscription := models.Subscription{
				PlanID:    newPlan.ID,
				UserID:    user.ID,
				StartDate: newStartDate,
				ValidTill: validTill,
			}

			tx.Create(&newSubscription)
		}
		return err
	})

	return res, err
}
