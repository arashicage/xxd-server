-- DROP TABLE IF EXISTS `im_chat`;
CREATE TABLE IF NOT EXISTS `im_chat` (
  `id` mediumint(8) unsigned NOT NULL AUTO_INCREMENT,
  `gid` char(40) NOT NULL DEFAULT '',
  `name` varchar(60) NOT NULL DEFAULT '',
  `type` varchar(20) NOT NULL DEFAULT 'group',
  `admins` varchar(255) NOT NULL DEFAULT '',
  `committers` varchar(255) NOT NULL DEFAULT '',
  `subject` mediumint(8) unsigned NOT NULL DEFAULT 0,
  `public` enum('0', '1') NOT NULL DEFAULT '0',
  `createdBy` varchar(30) NOT NULL DEFAULT '',
  `createdDate` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `editedBy` varchar(30) NOT NULL DEFAULT '',
  `editedDate` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `lastActiveTime` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`id`),
  KEY `gid` (`gid`),
  KEY `name` (`name`),
  KEY `type` (`type`),
  KEY `public` (`public`),
  KEY `createdBy` (`createdBy`),
  KEY `editedBy` (`editedBy`),
  UNIQUE KEY `im_chat` (`gid`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

-- DROP TABLE IF EXISTS `im_message`;
CREATE TABLE IF NOT EXISTS `im_message` (
  `id` mediumint(8) unsigned NOT NULL AUTO_INCREMENT,
  `gid` char(40) NOT NULL DEFAULT '',
  `cgid` char(40) NOT NULL DEFAULT '',
  `user` varchar(30) NOT NULL DEFAULT '',
  `date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `type` enum('normal', 'broadcast') NOT NULL DEFAULT 'normal',
  `content` text NOT NULL DEFAULT '',
  `contentType` enum('text', 'emotion', 'image', 'file', 'object') NOT NULL DEFAULT 'text',
  PRIMARY KEY (`id`),
  KEY `mgid` (`gid`),
  KEY `mcgid` (`cgid`),
  KEY `muser` (`user`),
  KEY `mtype` (`type`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

-- DROP TABLE IF EXISTS `im_chatuser`;
CREATE TABLE IF NOT EXISTS `im_chatuser`(
  `id` mediumint(8) unsigned NOT NULL AUTO_INCREMENT,
  `cgid` char(40) NOT NULL DEFAULT '',
  `user` mediumint(8) NOT NULL DEFAULT 0,
  `order` smallint(5) NOT NULL DEFAULT 0,
  `star` enum('0', '1') NOT NULL DEFAULT '0',
  `hide` enum('0', '1') NOT NULL DEFAULT '0',
  `mute` enum('0', '1') NOT NULL DEFAULT '0',
  `join` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `quit` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`id`),
  KEY `cgid` (`cgid`),
  KEY `user` (`user`),
  KEY `order` (`order`),
  KEY `star` (`star`),
  KEY `hide` (`hide`),
  UNIQUE KEY `chatuser` (`cgid`, `user`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

-- DROP TABLE IF EXISTS `im_usermessage`;
CREATE TABLE IF NOT EXISTS `im_usermessage`(
  `id` mediumint(8) NOT NULL AUTO_INCREMENT,
  `level` smallint(5) NOT NULL DEFAULT 3,
  `user` mediumint(8) NOT NULL DEFAULT 0,
  `message` text NOT NULL DEFAULT '',
  PRIMARY KEY `id` (`id`),
  KEY `muser` (`user`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `sys_user` (
	`id` MEDIUMINT(8) UNSIGNED NOT NULL AUTO_INCREMENT,
	`dept` MEDIUMINT(8) UNSIGNED NOT NULL,
	`account` CHAR(30) NOT NULL DEFAULT '',
	`password` CHAR(32) NOT NULL DEFAULT '',
	`realname` CHAR(30) NOT NULL DEFAULT '',
	`role` CHAR(30) NOT NULL,
	`nickname` CHAR(60) NOT NULL DEFAULT '',
	`admin` ENUM('no','common','super') NOT NULL DEFAULT 'no',
	`avatar` VARCHAR(255) NOT NULL DEFAULT '',
	`birthday` DATE NOT NULL,
	`gender` ENUM('f','m','u') NOT NULL DEFAULT 'u',
	`email` CHAR(90) NOT NULL DEFAULT '',
	`skype` CHAR(90) NOT NULL,
	`qq` CHAR(20) NOT NULL DEFAULT '',
	`yahoo` CHAR(90) NOT NULL DEFAULT '',
	`gtalk` CHAR(90) NOT NULL DEFAULT '',
	`wangwang` CHAR(90) NOT NULL DEFAULT '',
	`site` VARCHAR(100) NOT NULL,
	`mobile` CHAR(11) NOT NULL DEFAULT '',
	`phone` CHAR(20) NOT NULL DEFAULT '',
	`address` CHAR(120) NOT NULL DEFAULT '',
	`zipcode` CHAR(10) NOT NULL DEFAULT '',
	`visits` MEDIUMINT(8) UNSIGNED NOT NULL DEFAULT '0',
	`ip` CHAR(50) NOT NULL DEFAULT '',
	`last` DATETIME NOT NULL,
	`ping` DATETIME NOT NULL,
	`fails` TINYINT(3) UNSIGNED NOT NULL DEFAULT '0',
	`join` DATETIME NOT NULL,
	`locked` DATETIME NOT NULL,
	`deleted` ENUM('0','1') NOT NULL,
	`status` ENUM('online','away','busy','offline') NOT NULL DEFAULT 'offline',
	PRIMARY KEY (`id`),
	UNIQUE INDEX `account` (`account`),
	INDEX `admin` (`admin`),
	INDEX `accountPassword` (`account`, `password`),
	INDEX `dept` (`dept`)
)
COLLATE='utf8_general_ci'
ENGINE=MyISAM
AUTO_INCREMENT=6
;

CREATE TABLE IF NOT EXISTS `sys_file` (
	`id` MEDIUMINT(8) UNSIGNED NOT NULL AUTO_INCREMENT,
	`pathname` CHAR(100) NOT NULL,
	`title` CHAR(90) NOT NULL,
	`extension` CHAR(30) NOT NULL,
	`size` MEDIUMINT(8) UNSIGNED NOT NULL DEFAULT '0',
	`objectType` CHAR(30) NOT NULL,
	`objectID` MEDIUMINT(8) UNSIGNED NOT NULL,
	`createdBy` CHAR(30) NOT NULL DEFAULT '',
	`createdDate` DATETIME NOT NULL,
	`editor` ENUM('1','0') NOT NULL DEFAULT '0',
	`primary` ENUM('1','0') NULL DEFAULT '0',
	`public` ENUM('1','0') NOT NULL DEFAULT '1',
	`downloads` MEDIUMINT(8) UNSIGNED NOT NULL DEFAULT '0',
	`extra` VARCHAR(255) NOT NULL,
	PRIMARY KEY (`id`),
	INDEX `object` (`objectType`, `objectID`)
)
COLLATE='utf8_general_ci'
ENGINE=MyISAM
AUTO_INCREMENT=3
;

INSERT INTO `sys_user` (`id`, `dept`, `account`, `password`, `realname`, `role`, `nickname`, `admin`, `avatar`, `birthday`, `gender`, `email`, `skype`, `qq`, `yahoo`, `gtalk`, `wangwang`, `site`, `mobile`, `phone`, `address`, `zipcode`, `visits`, `ip`, `last`, `ping`, `fails`, `join`, `locked`, `deleted`, `status`) VALUES
	(1, 0, 'admin', 'c7122a1349c22cb3c009da3613d242ab', 'admin', '', '', 'super', '', '0000-00-00', 'u', '', '', '', '', '', '', '', '', '', '', '', 20, '172.30.11.230', '2017-11-09 13:24:05', '2017-11-09 13:24:05', 1, '2017-08-29 16:06:27', '0000-00-00 00:00:00', '0', 'offline'),


-- 1.3 update
ALTER TABLE `im_message` ADD `order` mediumint(8) unsigned NOT NULL AFTER `date`;
ALTER TABLE `im_message` ADD `data` text NOT NULL DEFAULT '' AFTER `contentType`;
ALTER TABLE `im_chatuser` ADD `category` varchar(40) NOT NULL DEFAULT '' AFTER `quit`;
ALTER TABLE `im_chat` ADD `dismissDate` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' AFTER `lastActiveTime`;

