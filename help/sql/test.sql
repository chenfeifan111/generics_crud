/*
 Navicat Premium Dump SQL

 Source Server         : local_mysql
 Source Server Type    : MySQL
 Source Server Version : 80012 (8.0.12)
 Source Host           : localhost:3306
 Source Schema         : test

 Target Server Type    : MySQL
 Target Server Version : 80012 (8.0.12)
 File Encoding         : 65001

 Date: 20/12/2025 15:53:45
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for group_example
-- ----------------------------
DROP TABLE IF EXISTS `group_example`;
CREATE TABLE `group_example`  (
  `id` int(11) NOT NULL,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `department` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '所在部门',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of group_example
-- ----------------------------
INSERT INTO `group_example` VALUES (1, 'zs', '开发部');
INSERT INTO `group_example` VALUES (2, 'ls', '开发部');
INSERT INTO `group_example` VALUES (3, 'ww', '运营部');
INSERT INTO `group_example` VALUES (4, 'zl', '销售部');
INSERT INTO `group_example` VALUES (5, 'tq', '运营部');

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `age` int(11) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES ('1', 'zs', 22);
INSERT INTO `user` VALUES ('2', 'ls', 17);
INSERT INTO `user` VALUES ('3', 'zss', 19);
INSERT INTO `user` VALUES ('4', 'zsss', 22);
INSERT INTO `user` VALUES ('5', 'ls2', 90);

SET FOREIGN_KEY_CHECKS = 1;
