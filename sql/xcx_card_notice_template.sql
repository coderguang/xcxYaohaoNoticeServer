/*
Navicat MySQL Data Transfer

Source Server         : 47.107.177.155
Source Server Version : 50505
Source Host           : 47.107.177.155:3306
Source Database       : xcx_template

Target Server Type    : MYSQL
Target Server Version : 50505
File Encoding         : 65001

Date: 2019-08-29 11:55:51
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for xcx_card_notice_template
-- ----------------------------
DROP TABLE IF EXISTS `xcx_card_notice_template`;
CREATE TABLE `xcx_card_notice_template` (
  `token_id` varchar(512) NOT NULL,
  `name` varchar(512) NOT NULL,
  `card_type` int(11) NOT NULL,
  `title` varchar(128) NOT NULL,
  `code` varchar(128) NOT NULL DEFAULT '' COMMENT '编码',
  `phone` varchar(128) NOT NULL DEFAULT '',
  `end_dt` datetime NOT NULL COMMENT '通知结束时间',
  `tips` varchar(1024) NOT NULL DEFAULT '',
  `renew_times` int(11) NOT NULL DEFAULT 0 COMMENT '续约次数，以token为key',
  `status` varchar(11) NOT NULL DEFAULT '0' COMMENT '0：正常   1：取消',
  `notice_times` int(11) NOT NULL DEFAULT 0 COMMENT '总通知次数',
  PRIMARY KEY (`token_id`,`title`),
  KEY `token` (`token_id`) USING BTREE,
  KEY `title_token` (`token_id`,`title`) USING BTREE,
  KEY `phone` (`phone`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
