-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `dietary_restrictions` (
  `dietary_restriction_id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `dietary_restriction_name` varchar(50) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`dietary_restriction_id`),
  CONSTRAINT `fk_user_dietary_restriction` FOREIGN KEY (`user_id`) REFERENCES `users`(`user_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `dietary_restrictions`;
-- +goose StatementEnd
