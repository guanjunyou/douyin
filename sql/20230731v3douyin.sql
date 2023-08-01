/*
 Navicat Premium Data Transfer

 Source Server         : mydp
 Source Server Type    : MySQL
 Source Server Version : 50714 (5.7.14)
 Source Host           : localhost:3306
 Source Schema         : douyin

 Target Server Type    : MySQL
 Target Server Version : 50714 (5.7.14)
 File Encoding         : 65001

 Date: 28/07/2023 17:30:49
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for comment
-- ----------------------------
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment`
(
    `id`          bigint(64)                                                   NOT NULL,
    `user_id`     bigint(64)                                                   NOT NULL COMMENT '评论用户的id',
    `content`     text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci        NOT NULL COMMENT '评论内容',
    `video_id`    varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '评论的视频id',
    `create_date` timestamp                                                    NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    `is_deleted`  int(1)                                                       NULL DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of comment
-- ----------------------------

-- ----------------------------
-- Table structure for follow
-- ----------------------------
DROP TABLE IF EXISTS `follow`;
CREATE TABLE `follow`
(
    `id`             bigint(64) NOT NULL AUTO_INCREMENT,
    `user_id`        bigint(64) NULL DEFAULT NULL COMMENT '用户id',
    `follow_user_id` bigint(64) NULL DEFAULT NULL COMMENT '关注的用户id',
    `create_date`    timestamp  NULL DEFAULT CURRENT_TIMESTAMP,
    `is_deleted`     int(1)     NULL DEFAULT 0,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of follow
-- ----------------------------

-- ----------------------------
-- Table structure for like
-- ----------------------------
DROP TABLE IF EXISTS `like`;
CREATE TABLE `like`
(
    `id`          bigint(64) NOT NULL,
    `video_id`    bigint(64) NULL DEFAULT NULL,
    `user_id`     bigint(64) NULL DEFAULT NULL,
    `create_date` timestamp  NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    `is_deleted`  int(1)     NULL DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of like
-- ----------------------------

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message`
(
    `id`          bigint(64)                                            NOT NULL,
    `content`     text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '消息内容',
    `create_date` timestamp                                             NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    `is_deleted`  int(1)                                                NULL DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of message
-- ----------------------------

-- ----------------------------
-- Table structure for message_push_event
-- ----------------------------
DROP TABLE IF EXISTS `message_push_event`;
CREATE TABLE `message_push_event`
(
    `id`           bigint(64)                                            NOT NULL,
    `from_user_id` bigint(64)                                            NULL DEFAULT NULL COMMENT '发送者的id',
    `msg_content`  text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '消息内容',
    `create_date`  timestamp                                             NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    `is_deleted`   int(1)                                                NULL DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of message_push_event
-- ----------------------------

-- ----------------------------
-- Table structure for message_send_event
-- ----------------------------
DROP TABLE IF EXISTS `message_send_event`;
CREATE TABLE `message_send_event`
(
    `id`          bigint(64)                                            NOT NULL,
    `user_id`     bigint(64)                                            NOT NULL,
    `to_user_id`  bigint(64)                                            NOT NULL,
    `msg_content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `create_date` timestamp                                             NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    `is_deleted`  int(1)                                                NULL DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of message_send_event
-- ----------------------------

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    `id`               bigint(64)                                                    NOT NULL COMMENT '用户id',
    `name`             varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '姓名',
    `follow_count`     int(8)                                                        NULL DEFAULT NULL COMMENT '关注数',
    `follower_count`   int(8)                                                        NULL DEFAULT NULL COMMENT '粉丝数',
    `phone`            varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL COMMENT '电话',
    `password`         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '密码',
    `avatar`           varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '头像',
    `gender`           int(2)                                                        NULL DEFAULT NULL COMMENT '性别',
    `age`              int(2)                                                        NULL DEFAULT NULL COMMENT '年龄',
    `create_date`      timestamp                                                     NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    `is_deleted`       int(1)                                                        NULL DEFAULT NULL,
    `nickname`         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '昵称',
    `signature`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '个人简介',
    `total_favorited`  int(22)                                                       NULL DEFAULT NULL COMMENT '获赞数量',
    `work_count`       int(22)                                                       NULL DEFAULT NULL COMMENT '作品数',
    `favorite_count`   int(22)                                                       NULL DEFAULT NULL COMMENT '喜欢数',
    `is_follow`        int(11)                                                       NULL DEFAULT NULL COMMENT '是否关注',
    `background_image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '个人背景图片',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user`
VALUES (7089783222816474111, '张四', 1, 1, '1', '1', '1', 1, 1, '2023-07-27 20:41:55', 0, '1', '0', 0, 0, 0, 1, NULL);
INSERT INTO `user`
VALUES (7090306410939941888, '20202231014@163.com', 0, 0, '',
        '$2a$10$t7RCzWVc1A/ReQPi8awWsu0MnnhAdwBTLzCsW1CWaHw1TU/64XIkG', '', 0, 0, '2023-07-27 20:26:20', 0, '', '', 0,
        0, 0, 0, '');

-- ----------------------------
-- Table structure for video
-- ----------------------------
DROP TABLE IF EXISTS `video`;
CREATE TABLE `video`
(
    `id`             bigint(64)                                                     NOT NULL,
    `author_id`      bigint(64)                                                     NOT NULL COMMENT '视频作者',
    `play_url`       varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '播放路径',
    `cover_url`      varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `favorite_count` int(8)                                                         NULL DEFAULT NULL COMMENT '喜欢数量',
    `comment_count`  int(8)                                                         NULL DEFAULT NULL COMMENT '评论数量',
    `is_favorite`    int(2)                                                         NULL DEFAULT NULL,
    `create_date`    timestamp                                                      NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    `is_deleted`     int(1)                                                         NULL DEFAULT NULL,
    `title`          varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NULL DEFAULT NULL COMMENT '视频标题',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of video
-- ----------------------------
INSERT INTO `video`
VALUES (7089783222816474111, 7089783222816474111, 'https://www.w3schools.com/html/movie.mp4',
        'https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg', 0, 0, 0, '2023-07-28 12:33:04', 0,
        '台风');
INSERT INTO `video`
VALUES (7089783222816474112, 7089783222816474111, 'https://cccimg.com/view.php/686384315ac21f0f67170063e07b1f75.mp4',
        'https://img1.imgtp.com/2023/07/28/PGnC0crf.png', 0, 0, 0, '2023-07-28 10:36:42', 0, '熊');

-- ----------------------------
-- Procedure structure for addFollowRelation
-- ----------------------------
DROP PROCEDURE IF EXISTS `addFollowRelation`;
delimiter ;;
CREATE PROCEDURE `addFollowRelation`(IN user_id bigint, IN follower_id bigint)
BEGIN
    #Routine body goes here...
    # 声明记录个数变量。
    DECLARE cnt INT DEFAULT 0;
    # 获取记录个数变量。
    SELECT COUNT(1) FROM follow f where f.user_id = user_id AND f.follow_user_id = follower_id INTO cnt;
    # 判断是否已经存在该记录，并做出相应的插入关系、更新关系动作。
    # 插入操作。
    IF cnt = 0 THEN
        INSERT INTO follow(`user_id`, `follow_user_id`) VALUES (user_id, follower_id);
    END IF;
    # 更新操作
    IF cnt != 0 THEN
        UPDATE follow f SET f.is_deleted = 0 WHERE f.user_id = user_id AND f.follow_user_id = follower_id;
    END IF;
END
;;
delimiter ;

-- ----------------------------
-- Procedure structure for delFollowRelation
-- ----------------------------
DROP PROCEDURE IF EXISTS `delFollowRelation`;
delimiter ;;
CREATE PROCEDURE `delFollowRelation`(IN `user_id` bigint, IN `follower_id` bigint)
BEGIN
    #Routine body goes here...
    # 定义记录个数变量，记录是否存在此关系，默认没有关系。
    DECLARE cnt INT DEFAULT 0;
    # 查看是否之前有关系。
    SELECT COUNT(1) FROM follow f WHERE f.user_id = user_id AND f.follow_user_id = follower_id INTO cnt;
# 有关系，则需要update cancel = 1，使其关系无效。
    IF cnt = 1 THEN
        UPDATE follow f SET f.is_deleted = 1 WHERE f.user_id = user_id AND f.follow_user_id = follower_id;
    END IF;
END
;;
delimiter ;

SET FOREIGN_KEY_CHECKS = 1;
