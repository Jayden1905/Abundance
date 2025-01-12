package types

type Goal struct {
	GoalID   int    `json:"goal_id"`
	UserID   int    `json:"user_id"`
	GoalName string `json:"goal_name"`
}

type CreateGoalsPayload struct {
	UserID string   `json:"user_id" validate:"required"`
	Goals  []string `json:"goals"`
}

type UpdateGoalsPayload struct {
	UserID string   `json:"user_id" validate:"required"`
	Goals  []string `json:"goals"`
}
