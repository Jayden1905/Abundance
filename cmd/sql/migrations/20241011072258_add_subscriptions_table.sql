-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `subscriptions` (
  `subscription_id` tinyint NOT NULL,
  `subscription_type` enum('Active', 'Inactive', 'Pending', 'Cancelled') NOT NULL,
  `start_date` date NOT NULL,
  `end_date` date NOT NULL,
  PRIMARY KEY (`subscription_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `subscriptions`;
-- +goose StatementEnd