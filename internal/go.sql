CREATE DATABASE golang;
USE golang;

# Dump of table bike_thefts
# ------------------------------------------------------------

DROP TABLE IF EXISTS `bike_thefts`;

CREATE TABLE `bike_thefts` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(45) NOT NULL DEFAULT '',
  `brand` varchar(45) NOT NULL DEFAULT '',
  `city` varchar(45) NOT NULL DEFAULT '',
  `description` longtext NOT NULL,
  `reported` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `solved` tinyint(1) NOT NULL DEFAULT '0',
  `officer` int(11) DEFAULT NULL,
  `image` blob,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



# Dump of table officers
# ------------------------------------------------------------

DROP TABLE IF EXISTS `officers`;

CREATE TABLE `officers` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

