SET foreign_key_checks=0;

DROP TABLE IF EXISTS `player`;

CREATE TABLE `player` (
    `id` BIGINT unsigned NOT NULL,
    `name` VARCHAR(191) NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    `daily_notification_at` TIME NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;


DROP TABLE IF EXISTS `entry`;

CREATE TABLE `entry` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `title` VARCHAR(100) NOT NULL,
    `public` TINYINT(1) NOT NULL DEFAULT 0,
    `content` TEXT NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    FULLTEXT `full_text_idx` (`content`) WITH PARSER ngram,
    INDEX `created_at_idx` (`created_at`),
    INDEX `title_idx` (`title`),
    UNIQUE `created_at_uniq_idx` (`created_at`),
    PRIMARY KEY (`id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;


DROP TABLE IF EXISTS `player_comment`;

CREATE TABLE `player_comment` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `player_id` INTEGER NOT NULL,
    `entry_id` INTEGER NOT NULL,
    `comment` VARCHAR(99) NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    INDEX `player_id_entry_id_idx` (`player_id`, `entry_id`),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;


DROP TABLE IF EXISTS `bookmark`;

CREATE TABLE `bookmark` (
    `id` INTEGER NOT NULL,
    `user_id` INTEGER NOT NULL,
    `entry_id` INTEGER NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    UNIQUE `user_id_entry_id` (`user_id`, `entry_id`),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET utf8mb4;

SET foreign_key_checks=1;