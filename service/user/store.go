package user

import (
	"database/sql"
	"fmt"

	"golang.org/x/net/context"

	"github.com/jayden1905/abundance/cmd/pkg/database"
	"github.com/jayden1905/abundance/types"
)

type Store struct {
	db *database.Queries
}

// NewStore initializes the Store with the database queries
func NewStore(db *database.Queries) *Store {
	return &Store{db: db}
}

// GetUsersPaginated fetches users by page from the database
func (s *Store) GetUsersPaginated(page int32, pageSize int32) ([]*types.User, error) {
	offset := (page - 1) * pageSize
	users, err := s.db.GetAllUsersPaginated(context.Background(), database.GetAllUsersPaginatedParams{
		Limit:  pageSize,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	// Convert the database user to the user type
	var allUsers []*types.User

	for _, user := range users {
		allUsers = append(allUsers, &types.User{
			ID:           int32(user.UserID),
			Username:     user.Username,
			Role:         string(user.Role),
			Subscription: string(user.SubscriptionType),
			Email:        user.Email,
			IsVerified:   user.IsVerified,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
		})
	}

	return allUsers, nil
}

// GetUserByEmail fetches a user by email from the database
func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	user, err := s.db.GetUserByEmail(context.Background(), email) // Use the SQLC-generated method
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &types.User{
		ID:           int32(user.UserID),
		Username:     user.Username,
		Role:         string(user.Role),
		Subscription: string(user.SubscriptionType),
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		IsVerified:   user.IsVerified,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

// GetUserByID fetches a user by ID from the database
func (s *Store) GetUserByID(id int32) (*types.User, error) {
	user, err := s.db.GetUserByID(context.Background(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &types.User{
		ID:           int32(user.UserID),
		Username:     user.Username,
		Role:         string(user.Role),
		Subscription: string(user.SubscriptionType),
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		IsVerified:   user.IsVerified,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

// GetUserRoleByID fetches the role of a user by ID from the database
func (s *Store) GetUserRoleByID(id int32) (string, error) {
	role, err := s.db.GetUserRoleByUserID(context.Background(), id)
	if err != nil {
		return "", err
	}

	stringRole := string(role)

	return stringRole, nil
}

// CreateUser creates a new user in the database
func (s *Store) CreateUser(ctx context.Context, user *database.User) error {
	err := s.db.CreateUser(ctx, database.CreateUserParams{
		Username:       user.Username,
		Email:          user.Email,
		PasswordHash:   user.PasswordHash,
		RoleID:         user.RoleID,
		SubscriptionID: user.SubscriptionID,
	})
	if err != nil {
		return err
	}

	return nil
}

// CreateSuperUser creates a new super user in the database
func (s *Store) CreateSuperUser(ctx context.Context, user *types.User) error {
	err := s.db.CreateAdmin(ctx, database.CreateAdminParams{
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	})
	if err != nil {
		return err
	}

	return nil
}

// Update the user to a super user
func (s *Store) UpdateUserToSuperUser(ctx context.Context, id int32) error {
	err := s.db.UpdateUserToAdmin(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// Update the user to a normal user
func (s *Store) UpdateUserToNormalUser(ctx context.Context, id int32) error {
	err := s.db.UpdateUserToFreeUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUserVerification updates the user verification status in the database
func (s *Store) UpdateUserVerification(ctx context.Context, id int32) error {
	err := s.db.UpdateUserVerificationStatus(ctx, database.UpdateUserVerificationStatusParams{
		IsVerified: true,
		UserID:     id,
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteUserByID deletes a user by ID from the database
func (s *Store) DeleteUserByID(ctx context.Context, id int32) error {
	err := s.db.DeleteUserByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUserPassword updates the user password in the database
func (s *Store) UpdateUserPassword(ctx context.Context, id int32, passwordHash string) error {
	err := s.db.UpdateUserPassword(ctx, database.UpdateUserPasswordParams{
		PasswordHash: passwordHash,
		UserID:       id,
	})
	if err != nil {
		return err
	}

	return nil
}
