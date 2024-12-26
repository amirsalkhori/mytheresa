CREATE TABLE `prices` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `original` INT UNSIGNED NOT NULL,
    `final` INT UNSIGNED NOT NULL,
    `discount_percentage` TINYINT UNSIGNED NOT NULL,
    `currency` VARCHAR(10) NOT NULL,
    INDEX `idx_currency` (`currency`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
