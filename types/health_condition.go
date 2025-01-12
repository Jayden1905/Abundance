package types

type HealthCondition struct {
	HealthConditionID int    `json:"health_condition_id"`
	UserID            int    `json:"user_id"`
	ConditionName     string `json:"condition_name"`
}

type CreateHealthConditionPayload struct {
	UserID     string   `json:"user_id" validate:"required"`
	Conditions []string `json:"conditions"`
}

type UpdateHealthConditionPayload struct {
	UserID     string   `json:"user_id" validate:"required"`
	Conditions []string `json:"conditions"`
}
