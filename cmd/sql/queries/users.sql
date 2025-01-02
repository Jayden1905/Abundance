-- name: CreateUser :exec
INSERT INTO users (
        role_id,
        username,
        email,
        password_hash,
        subscription_id
    )
VALUES (?, ?, ?, ?, ?);
-- name: CreateAdmin :exec
INSERT INTO users (
        role_id,
        username,
        email,
        password_hash,
        subscription_id
    )
VALUES (4, ?, ?, ?, 1);
-- name: GetAllUsersPaginated :many
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
LIMIT ? OFFSET ?;
-- name: GetUserByID :one
SELECT users.user_id,
    roles.name AS 'role',
    users.username,
    users.email,
    users.password_hash,
    users.is_verified,
    subscriptions.subscription_type AS 'subscription type',
    users.created_at,
    users.updated_at
FROM users users
    JOIN roles roles USING(role_id)
    JOIN subscriptions subscriptions USING (subscription_id)
WHERE user_id = ?;
-- name: GetUserByEmail :one
SELECT users.user_id,
    roles.name AS 'role',
    users.username,
    users.email,
    subscriptions.subscription_type AS 'subscription type',
    users.is_verified,
    users.created_at,
    users.updated_at,
    users.password_hash
FROM users users
    JOIN roles roles USING(role_id)
    JOIN subscriptions subscriptions USING (subscription_id)
WHERE email = ?;
-- name: UpdateUserToAdmin :exec
UPDATE users
SET role_id = 4
WHERE user_id = ?;
-- name: GetUserRoleByUserID :one
SELECT roles.name
FROM users users
    JOIN roles roles using(role_id)
WHERE user_id = ?;
-- name: UpdateUserToFreeUser :exec
UPDATE users
SET role_id = 1
WHERE user_id = ?;
-- name: UpdateUserSubscriptionStatus :exec
UPDATE users
SET subscription_id = ?
WHERE user_id = ?;
-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = ?
WHERE user_id = ?;
-- name: DeleteUserByID :exec
DELETE FROM users
WHERE user_id = ?;
-- name: UpdateUserVerificationStatus :exec
UPDATE users
SET is_verified = ?
WHERE user_id = ?;
