-- +goose Up
-- Seed roles and subscriptions individually
INSERT IGNORE INTO roles (role_id, name)
VALUES (1, 'free_user');
INSERT IGNORE INTO roles (role_id, name)
VALUES (2, 'premium_user');
INSERT IGNORE INTO roles (role_id, name)
VALUES (3, 'nutritionist');
INSERT IGNORE INTO roles (role_id, name)
VALUES (4, 'admin');
INSERT IGNORE INTO subscriptions (subscription_id, subscription_type)
VALUES (1, 'Active');
INSERT IGNORE INTO subscriptions (subscription_id, subscription_type)
VALUES (2, 'Inactive');
INSERT IGNORE INTO subscriptions (subscription_id, subscription_type)
VALUES (3, 'Pending');
INSERT IGNORE INTO subscriptions (subscription_id, subscription_type)
VALUES (4, 'Cancelled');
-- +goose Down
-- Rollback seed roles and subscriptions
DELETE FROM roles
WHERE role_id IN (1, 2);
DELETE FROM subscriptions
WHERE subscription_id IN (1, 2, 3, 4);