CREATE TABLE `accounts` (
  `id` bigint AUTO_INCREMENT PRIMARY KEY,
  `owner` varchar(255) NOT NULL,
  `balance` bigint NOT NULL,
  `currency` varchar(255) NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `entries` (
  `id` bigint AUTO_INCREMENT PRIMARY KEY,
  `account_id` bigint,
  `amount` bigint NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`)
);

CREATE TABLE `transfers` (
  `id` bigint AUTO_INCREMENT PRIMARY KEY,
  `from_account` bigint,
  `to_account` bigint,
  `amount` bigint NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`from_account`) REFERENCES `accounts` (`id`),
  FOREIGN KEY (`to_account`) REFERENCES `accounts` (`id`)
);

CREATE INDEX `accounts_index_0` ON `accounts` (`owner`);
CREATE INDEX `entries_index_1` ON `entries` (`account_id`);
CREATE INDEX `transfers_index_2` ON `transfers` (`from_account`);
CREATE INDEX `transfers_index_3` ON `transfers` (`to_account`);
CREATE INDEX `transfers_index_4` ON `transfers` (`from_account`, `to_account`);
