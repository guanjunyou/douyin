/*
Navicat MySQL Data Transfer

Source Server         : 本地主机
Source Server Version : 50724
Source Host           : localhost:3306
Source Database       : douyin

Target Server Type    : MYSQL
Target Server Version : 50724
File Encoding         : 65001

Date: 2023-07-26 15:09:13
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for comment
-- ----------------------------
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment` (
  `id` bigint(64) NOT NULL,
  `user_id` bigint(64) NOT NULL COMMENT '评论用户的id',
  `content` text NOT NULL COMMENT '评论内容',
  `vedio_id` varchar(64) NOT NULL COMMENT '评论的视频id',
  `create_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `is_deleted` int(1) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of comment
-- ----------------------------

-- ----------------------------
-- Table structure for follow
-- ----------------------------
DROP TABLE IF EXISTS `follow`;
CREATE TABLE `follow` (
  `id` bigint(64) NOT NULL,
  `user_id` bigint(64) DEFAULT NULL COMMENT '用户id',
  `follow_user_id` bigint(64) DEFAULT NULL COMMENT '关注的用户id',
  `create_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `is_deleted` int(1) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of follow
-- ----------------------------

-- ----------------------------
-- Table structure for like
-- ----------------------------
DROP TABLE IF EXISTS `like`;
CREATE TABLE `like` (
  `id` bigint(64) NOT NULL,
  `vedio_id` bigint(64) DEFAULT NULL,
  `user_id` bigint(64) DEFAULT NULL,
  `create_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `is_deleted` int(1) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of like
-- ----------------------------

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message` (
  `id` bigint(64) NOT NULL,
  `content` text NOT NULL COMMENT '消息内容',
  `create_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `is_deleted` int(1) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of message
-- ----------------------------

-- ----------------------------
-- Table structure for message_push_event
-- ----------------------------
DROP TABLE IF EXISTS `message_push_event`;
CREATE TABLE `message_push_event` (
  `id` bigint(64) NOT NULL,
  `from_user_id` bigint(64) DEFAULT NULL COMMENT '发送者的id',
  `msg_content` text COMMENT '消息内容',
  `create_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `is_deleted` int(1) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of message_push_event
-- ----------------------------

-- ----------------------------
-- Table structure for message_send_event
-- ----------------------------
DROP TABLE IF EXISTS `message_send_event`;
CREATE TABLE `message_send_event` (
  `id` bigint(64) NOT NULL,
  `user_id` bigint(64) NOT NULL,
  `to_user_id` bigint(64) NOT NULL,
  `msg_content` text NOT NULL,
  `create_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `is_deleted` int(1) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of message_send_event
-- ----------------------------

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint(64) NOT NULL COMMENT '用户id',
  `username` varchar(20) DEFAULT NULL COMMENT '姓名',
  `follow_count` int(8) DEFAULT NULL COMMENT '关注数',
  `follower_count` int(8) DEFAULT NULL COMMENT '粉丝数',
  `phone` varchar(11) NOT NULL COMMENT '电话',
  `password` varchar(255) DEFAULT NULL COMMENT '密码',
  `icon` varchar(255) DEFAULT NULL COMMENT '头像',
  `gender` int(2) DEFAULT NULL COMMENT '性别',
  `age` int(2) DEFAULT NULL COMMENT '年龄',
  `create_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `is_deleted` int(1) DEFAULT NULL,
  `nickname` varchar(255) DEFAULT NULL COMMENT '昵称',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES ('7089783222816474111', '张三', '1', '1', '1', '1', '1', '1', '1', '2023-07-26 13:35:27', '0', '1');

-- ----------------------------
-- Table structure for video
-- ----------------------------
DROP TABLE IF EXISTS `video`;
CREATE TABLE `video` (
  `id` bigint(64) NOT NULL,
  `author_id` bigint(64) NOT NULL COMMENT '视频作者',
  `play_url` varchar(2048) NOT NULL COMMENT '播放路径',
  `cover_url` varchar(2048) NOT NULL,
  `favorite_count` int(8) DEFAULT NULL COMMENT '喜欢数量',
  `comment_count` int(8) DEFAULT NULL COMMENT '评论数量',
  `is_favorite` int(2) DEFAULT NULL,
  `create_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `is_deleted` int(1) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of video
-- ----------------------------
INSERT INTO `video` VALUES ('7089783222816474112', '7089783222816474111', 'https://www.w3schools.com/html/movie.mp4', 'https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg', '0', '0', '0', '2023-07-26 13:18:01', '0');

-- ----------------------------
-- Procedure structure for addFollowRelation
-- ----------------------------
DROP PROCEDURE IF EXISTS `addFollowRelation`;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `addFollowRelation`(IN user_id int,IN follower_id int)
BEGIN
	#Routine body goes here...
	# 声明记录个数变量。
	DECLARE cnt INT DEFAULT 0;
	# 获取记录个数变量。
	SELECT COUNT(1) FROM follows f where f.user_id = user_id AND f.follower_id = follower_id INTO cnt;
	# 判断是否已经存在该记录，并做出相应的插入关系、更新关系动作。
	# 插入操作。
	IF cnt = 0 THEN
		INSERT INTO follows(`user_id`,`follower_id`) VALUES(user_id,follower_id);
	END IF;
	# 更新操作
	IF cnt != 0 THEN
		UPDATE follows f SET f.cancel = 0 WHERE f.user_id = user_id AND f.follower_id = follower_id;
	END IF;
END
;;
DELIMITER ;

-- ----------------------------
-- Procedure structure for delFollowRelation
-- ----------------------------
DROP PROCEDURE IF EXISTS `delFollowRelation`;
DELIMITER ;;
CREATE DEFINER=`root`@`localhost` PROCEDURE `delFollowRelation`(IN `user_id` int,IN `follower_id` int)
BEGIN
	#Routine body goes here...
	# 定义记录个数变量，记录是否存在此关系，默认没有关系。
	DECLARE cnt INT DEFAULT 0;
	# 查看是否之前有关系。
	SELECT COUNT(1) FROM follows f WHERE f.user_id = user_id AND f.follower_id = follower_id INTO cnt;
	# 有关系，则需要update cancel = 1，使其关系无效。
	IF cnt = 1 THEN
		UPDATE follows f SET f.cancel = 1 WHERE f.user_id = user_id AND f.follower_id = follower_id;
	END IF;
END
;;
DELIMITER ;
