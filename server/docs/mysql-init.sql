CREATE DATABASE IF NOT EXISTS authserver
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

GRANT ALL PRIVILEGES ON authserver.* TO 'auth'@'localhost' IDENTIFIED BY 'auth';
GRANT ALL PRIVILEGES ON authserver.* TO 'auth'@'127.0.0.1' IDENTIFIED BY 'auth';
GRANT ALL PRIVILEGES ON authserver.* TO 'auth'@'192.168.65.%' IDENTIFIED BY 'auth';
GRANT ALL PRIVILEGES ON authserver.* TO 'auth'@'%' IDENTIFIED BY 'auth';

FLUSH PRIVILEGES;

USE authserver;

CREATE TABLE IF NOT EXISTS `business_lines` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_business_lines_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `business_line_users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `business_line_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `permission` bigint NOT NULL DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_business_line_users_unique` (`business_line_id`,`user_id`),
  KEY `idx_business_line_users_business_line_id` (`business_line_id`),
  KEY `idx_business_line_users_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `business_line_wayne_namespaces` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `business_line_id` bigint unsigned NOT NULL,
  `wayne_namespace_id` bigint unsigned NOT NULL,
  `wayne_namespace_name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `kube_namespace` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_business_line_wayne_namespaces_unique` (`business_line_id`,`wayne_namespace_id`),
  KEY `idx_business_line_wayne_namespaces_business_line_id` (`business_line_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
