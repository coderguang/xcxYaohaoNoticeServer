/*
Navicat MySQL Data Transfer

Source Server         : 47.107.177.155
Source Server Version : 50505
Source Host           : 47.107.177.155:3306
Source Database       : xcx_template

Target Server Type    : MYSQL
Target Server Version : 50505
File Encoding         : 65001

Date: 2019-09-23 21:01:05
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for xcx_card_notice_require_data_template
-- ----------------------------
DROP TABLE IF EXISTS `xcx_card_notice_require_data_template`;
CREATE TABLE `xcx_card_notice_require_data_template` (
  `token_id` varchar(512) NOT NULL,
  `title` varchar(128) NOT NULL,
  `require_time` int(11) DEFAULT NULL,
  `name` varchar(512) NOT NULL,
  `final_login` datetime NOT NULL COMMENT '通知结束时间',
  `share_times` int(11) NOT NULL DEFAULT 0 COMMENT '续约次数，以token为key',
  `tips` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`token_id`,`title`),
  KEY `token` (`token_id`) USING BTREE,
  KEY `title_token` (`token_id`,`title`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
