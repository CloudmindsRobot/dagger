# ************************************************************
# Sequel Pro SQL dump
# Version 4541
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 10.12.32.179 (MySQL 5.6.46)
# Database: log
# Generation Time: 2021-01-20 01:54:06 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

CREATE DATABASE IF NOT EXISTS log DEFAULT CHARACTER SET = utf8;

# Dump of table auth_user
# ------------------------------------------------------------

DROP TABLE IF EXISTS `auth_user`;

CREATE TABLE `auth_user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `is_superuser` tinyint(4) NOT NULL DEFAULT '0',
  `is_active` tinyint(4) NOT NULL DEFAULT '1',
  `username` varchar(128) NOT NULL,
  `password` varchar(256) DEFAULT NULL,
  `mobile` varchar(32) DEFAULT NULL,
  `email` varchar(128) DEFAULT NULL,
  `create_at` datetime NOT NULL,
  `last_login_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table log_event
# ------------------------------------------------------------

DROP TABLE IF EXISTS `log_event`;

CREATE TABLE `log_event` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `resolve_at` datetime DEFAULT NULL,
  `create_at` datetime DEFAULT NULL,
  `rule_id` bigint(20) DEFAULT NULL,
  `status` varchar(24) NOT NULL,
  `count` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_log_event_rule_id` (`rule_id`),
  CONSTRAINT `fk_log_event_log_rule` FOREIGN KEY (`rule_id`) REFERENCES `log_rule` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table log_event_detail
# ------------------------------------------------------------

DROP TABLE IF EXISTS `log_event_detail`;

CREATE TABLE `log_event_detail` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `starts_at` datetime DEFAULT NULL,
  `summary` varchar(512) DEFAULT NULL,
  `labels` longtext,
  `description` longtext NOT NULL,
  `rule_id` bigint(20) DEFAULT NULL,
  `event_id` bigint(20) DEFAULT NULL,
  `level` varchar(32) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_log_event_detail_starts_at` (`starts_at`),
  KEY `idx_log_event_detail_rule_id` (`rule_id`),
  KEY `idx_log_event_detail_event_id` (`event_id`),
  CONSTRAINT `fk_log_event_detail_log_event` FOREIGN KEY (`event_id`) REFERENCES `log_event` (`id`),
  CONSTRAINT `fk_log_event_detail_log_rule` FOREIGN KEY (`rule_id`) REFERENCES `log_rule` (`id`),
  CONSTRAINT `fk_log_event_details` FOREIGN KEY (`event_id`) REFERENCES `log_event` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table log_group
# ------------------------------------------------------------

DROP TABLE IF EXISTS `log_group`;

CREATE TABLE `log_group` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `group_id` bigint(20) DEFAULT NULL,
  `rule_id` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_log_group_log_user_group_id` (`group_id`),
  KEY `idx_log_group_rule_id` (`rule_id`),
  CONSTRAINT `fk_log_group_log_rule` FOREIGN KEY (`rule_id`) REFERENCES `log_rule` (`id`),
  CONSTRAINT `fk_log_group_log_user_group` FOREIGN KEY (`group_id`) REFERENCES `log_user_group` (`id`),
  CONSTRAINT `fk_log_rule_groups` FOREIGN KEY (`rule_id`) REFERENCES `log_rule` (`id`),
  CONSTRAINT `fk_log_user_group_groups` FOREIGN KEY (`group_id`) REFERENCES `log_user_group` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table log_history
# ------------------------------------------------------------

DROP TABLE IF EXISTS `log_history`;

CREATE TABLE `log_history` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `label_json` longtext,
  `create_at` datetime DEFAULT NULL,
  `filter_json` longtext,
  `user_id` bigint(20) DEFAULT NULL,
  `log_ql` varchar(1024) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_log_history_user_id` (`user_id`),
  CONSTRAINT `fk_log_history_user` FOREIGN KEY (`user_id`) REFERENCES `auth_user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table log_label
# ------------------------------------------------------------

DROP TABLE IF EXISTS `log_label`;

CREATE TABLE `log_label` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `key` varchar(128) DEFAULT NULL,
  `value` varchar(128) DEFAULT NULL,
  `rule_id` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_log_label_rule_id` (`rule_id`),
  CONSTRAINT `fk_log_label_log_rule` FOREIGN KEY (`rule_id`) REFERENCES `log_rule` (`id`),
  CONSTRAINT `fk_log_rule_labels` FOREIGN KEY (`rule_id`) REFERENCES `log_rule` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table log_rule
# ------------------------------------------------------------

DROP TABLE IF EXISTS `log_rule`;

CREATE TABLE `log_rule` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `create_at` datetime DEFAULT NULL,
  `key` varchar(128) NOT NULL,
  `name` varchar(64) DEFAULT NULL,
  `description` varchar(2056) DEFAULT NULL,
  `summary` varchar(2056) DEFAULT NULL,
  `log_ql` varchar(512) DEFAULT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `update_at` datetime DEFAULT NULL,
  `level` varchar(32) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `key` (`key`),
  KEY `idx_log_rule_user_id` (`user_id`),
  CONSTRAINT `fk_log_rule_user` FOREIGN KEY (`user_id`) REFERENCES `auth_user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table log_snapshot
# ------------------------------------------------------------

DROP TABLE IF EXISTS `log_snapshot`;

CREATE TABLE `log_snapshot` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL,
  `count` int(11) DEFAULT NULL,
  `create_at` datetime DEFAULT NULL,
  `download_url` varchar(512) DEFAULT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `start_time` datetime DEFAULT NULL,
  `end_time` datetime DEFAULT NULL,
  `dir` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_log_snapshot_user_id` (`user_id`),
  CONSTRAINT `fk_log_snapshot_user` FOREIGN KEY (`user_id`) REFERENCES `auth_user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table log_user
# ------------------------------------------------------------

DROP TABLE IF EXISTS `log_user`;

CREATE TABLE `log_user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `group_id` bigint(20) DEFAULT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_log_user_log_user_group_id` (`group_id`),
  KEY `idx_log_user_user_id` (`user_id`),
  CONSTRAINT `fk_log_user_group_users` FOREIGN KEY (`group_id`) REFERENCES `log_user_group` (`id`),
  CONSTRAINT `fk_log_user_log_user_group` FOREIGN KEY (`group_id`) REFERENCES `log_user_group` (`id`),
  CONSTRAINT `fk_log_user_user` FOREIGN KEY (`user_id`) REFERENCES `auth_user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table log_user_group
# ------------------------------------------------------------

DROP TABLE IF EXISTS `log_user_group`;

CREATE TABLE `log_user_group` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `create_at` datetime DEFAULT NULL,
  `group_name` varchar(64) DEFAULT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_log_user_group_user_id` (`user_id`),
  CONSTRAINT `fk_log_user_group_user` FOREIGN KEY (`user_id`) REFERENCES `auth_user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;

USE log;

INSERT INTO `auth_user` (`is_superuser`, `is_active`, `username`, `password`, `mobile`, `email`, `create_at`, `last_login_at`)
VALUES (1, 1, 'admin',  '$2a$10$zqlCha8VIdeXeixuwFDlAerOFaimREojlZdDfqhPn3dwYbdD9T8n6', NULL, NULL, now(), now());
