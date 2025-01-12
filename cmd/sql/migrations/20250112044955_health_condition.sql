-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `health_conditions` (
  `health_condition_id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `health_condition_name` varchar(50) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`health_condition_id`),
  CONSTRAINT `fk_user_health_condition` FOREIGN KEY (`user_id`) REFERENCES `users`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE

) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `health_conditions`;
-- +goose StatementEnd
