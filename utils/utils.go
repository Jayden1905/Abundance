package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/jayden1905/abundance/types"
)

var Validate = validator.New()

func ValidatePayload(payload interface{}) (map[string]string, error) {
	err := Validate.Struct(payload)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			invalidFields := make(map[string]string)
			for _, e := range validationErrors {
				invalidFields[e.Field()] = fmt.Sprintf("Validation failed on the '%s' tag", e.Tag())
			}
			return invalidFields, fmt.Errorf("invalid payload")
		}
	}
	return nil, nil
}

func IsSuperUser(userID int32, store types.UserStore) (bool, error) {
	// Get the user role from the store
	role, err := store.GetUserRoleByID(userID)
	if err != nil {
		return false, fmt.Errorf("error getting user role by id: %v", err)
	}

	// Check if the role is "super_user"
	if role == "super_user" {
		return true, nil
	}

	return false, nil
}

func ConvertRoleStringToRoleID(role string) int8 {
	switch role {
	case "free_user":
		return 1
	case "premium_user":
		return 2
	case "nutritionist":
		return 3
	case "admin":
		return 4
	default:
		return 1
	}
}

func ConvertSubscriptionStringToSubscriptionID(subscription string) int8 {
	switch subscription {
	case "Active":
		return 1
	case "Inactive":
		return 2
	case "Pending":
		return 3
	case "Cancelled":
		return 4
	default:
		return 1
	}
}
