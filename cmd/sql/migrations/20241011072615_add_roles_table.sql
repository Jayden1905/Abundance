-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `roles` (
  `role_id` tinyint NOT NULL,
  `name` enum(
    'free_user',
    'premium_user',
    'nutritionist',
    'admin'
  ) NOT NULL,
  PRIMARY KEY (`role_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `roles`;
-- +goose StatementEnd