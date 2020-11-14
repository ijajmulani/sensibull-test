package subscriptions

type PostArgs struct {
	StartDate string `json:"start_date"`
	PlanName  string `json:"plan_id"`
	UserName  string `json:"user_name"`
}
