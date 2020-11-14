package subscriptions

type PostArgs struct {
	StartDate string `json:"start_date"`
	PlanID    string `json:"plan_id"`
	UserName  string `json:"user_name"`
}
