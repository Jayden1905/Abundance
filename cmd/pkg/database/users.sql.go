// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"
	"time"
)

const createAdmin = `-- name: CreateAdmin :exec
INSERT INTO users (
        role_id,
        username,
        email,
        password_hash,
        subscription_id
    )
VALUES (4, ?, ?, ?, 1)
`

type CreateAdminParams struct {
	Username     string
	Email        string
	PasswordHash string
}

func (q *Queries) CreateAdmin(ctx context.Context, arg CreateAdminParams) error {
	_, err := q.db.ExecContext(ctx, createAdmin, arg.Username, arg.Email, arg.PasswordHash)
	return err
}

const createUser = `-- name: CreateUser :exec
INSERT INTO users (
        role_id,
        username,
        email,
        password_hash,
        subscription_id
    )
VALUES (?, ?, ?, ?, ?)
`

type CreateUserParams struct {
	RoleID         int8
	Username       string
	Email          string
	PasswordHash   string
	SubscriptionID int8
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.ExecContext(ctx, createUser,
		arg.RoleID,
		arg.Username,
		arg.Email,
		arg.PasswordHash,
		arg.SubscriptionID,
	)
	return err
}

const deleteUserByID = `-- name: DeleteUserByID :exec
DELETE FROM users
WHERE user_id = ?
`

func (q *Queries) DeleteUserByID(ctx context.Context, userID int32) error {
	_, err := q.db.ExecContext(ctx, deleteUserByID, userID)
	return err
}

const getAllUsersPaginated = `-- name: GetAllUsersPaginated :many
SELECT users.user_id,
    roles.name AS 'role',
    users.username,
    users.email,
    users.is_verified,
    subscriptions.subscription_type AS 'subscription type',
    users.created_at,
    users.updated_at
FROM users
    JOIN roles USING(role_id)
    JOIN subscriptions USING (subscription_id)
ORDER BY users.created_at DESC
LIMIT ? OFFSET ?
`

type GetAllUsersPaginatedParams struct {
	Limit  int32
	Offset int32
}

type GetAllUsersPaginatedRow struct {
	UserID           int32
	Role             RolesName
	Username         string
	Email            string
	IsVerified       bool
	SubscriptionType SubscriptionsSubscriptionType
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (q *Queries) GetAllUsersPaginated(ctx context.Context, arg GetAllUsersPaginatedParams) ([]GetAllUsersPaginatedRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllUsersPaginated, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllUsersPaginatedRow
	for rows.Next() {
		var i GetAllUsersPaginatedRow
		if err := rows.Scan(
			&i.UserID,
			&i.Role,
			&i.Username,
			&i.Email,
			&i.IsVerified,
			&i.SubscriptionType,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT users.user_id,
    roles.name AS 'role',
    users.username,
    users.email,
    subscriptions.subscription_type AS 'subscription type',
    users.is_verified,
    users.created_at,
    users.updated_at
FROM users users
    JOIN roles roles USING(role_id)
    JOIN subscriptions subscriptions USING (subscription_id)
WHERE email = ?
`

type GetUserByEmailRow struct {
	UserID           int32
	Role             RolesName
	Username         string
	Email            string
	SubscriptionType SubscriptionsSubscriptionType
	IsVerified       bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i GetUserByEmailRow
	err := row.Scan(
		&i.UserID,
		&i.Role,
		&i.Username,
		&i.Email,
		&i.SubscriptionType,
		&i.IsVerified,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT users.user_id,
    roles.name AS 'role',
    users.username,
    users.email,
    users.is_verified,
    subscriptions.subscription_type AS 'subscription type',
    users.created_at,
    users.updated_at
FROM users users
    JOIN roles roles USING(role_id)
    JOIN subscriptions subscriptions USING (subscription_id)
WHERE user_id = ?
`

type GetUserByIDRow struct {
	UserID           int32
	Role             RolesName
	Username         string
	Email            string
	IsVerified       bool
	SubscriptionType SubscriptionsSubscriptionType
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (q *Queries) GetUserByID(ctx context.Context, userID int32) (GetUserByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, userID)
	var i GetUserByIDRow
	err := row.Scan(
		&i.UserID,
		&i.Role,
		&i.Username,
		&i.Email,
		&i.IsVerified,
		&i.SubscriptionType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserRoleByUserID = `-- name: GetUserRoleByUserID :one
SELECT roles.name
FROM users users
    JOIN roles roles using(role_id)
WHERE user_id = ?
`

func (q *Queries) GetUserRoleByUserID(ctx context.Context, userID int32) (RolesName, error) {
	row := q.db.QueryRowContext(ctx, getUserRoleByUserID, userID)
	var name RolesName
	err := row.Scan(&name)
	return name, err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = ?
WHERE user_id = ?
`

type UpdateUserPasswordParams struct {
	PasswordHash string
	UserID       int32
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.ExecContext(ctx, updateUserPassword, arg.PasswordHash, arg.UserID)
	return err
}

const updateUserSubscriptionStatus = `-- name: UpdateUserSubscriptionStatus :exec
UPDATE users
SET subscription_id = ?
WHERE user_id = ?
`

type UpdateUserSubscriptionStatusParams struct {
	SubscriptionID int8
	UserID         int32
}

func (q *Queries) UpdateUserSubscriptionStatus(ctx context.Context, arg UpdateUserSubscriptionStatusParams) error {
	_, err := q.db.ExecContext(ctx, updateUserSubscriptionStatus, arg.SubscriptionID, arg.UserID)
	return err
}

const updateUserToAdmin = `-- name: UpdateUserToAdmin :exec
UPDATE users
SET role_id = 4
WHERE user_id = ?
`

func (q *Queries) UpdateUserToAdmin(ctx context.Context, userID int32) error {
	_, err := q.db.ExecContext(ctx, updateUserToAdmin, userID)
	return err
}

const updateUserToFreeUser = `-- name: UpdateUserToFreeUser :exec
UPDATE users
SET role_id = 1
WHERE user_id = ?
`

func (q *Queries) UpdateUserToFreeUser(ctx context.Context, userID int32) error {
	_, err := q.db.ExecContext(ctx, updateUserToFreeUser, userID)
	return err
}

const updateUserVerificationStatus = `-- name: UpdateUserVerificationStatus :exec
UPDATE users
SET is_verified = ?
WHERE user_id = ?
`

type UpdateUserVerificationStatusParams struct {
	IsVerified bool
	UserID     int32
}

func (q *Queries) UpdateUserVerificationStatus(ctx context.Context, arg UpdateUserVerificationStatusParams) error {
	_, err := q.db.ExecContext(ctx, updateUserVerificationStatus, arg.IsVerified, arg.UserID)
	return err
}