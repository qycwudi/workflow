-- Modify "api" table
ALTER TABLE `api` ADD COLUMN `history_id` int NOT NULL;
-- Modify "api_secret_key" table
ALTER TABLE `api_secret_key` ADD INDEX `nidx_secret_key` (`secret_key`);
-- Modify "datasource" table
ALTER TABLE `datasource` ADD UNIQUE INDEX `un_name` (`name`) COMMENT "名称唯一索引";
-- Modify "trace" table
ALTER TABLE `trace` MODIFY COLUMN `workspace_id` varchar(255) NOT NULL COMMENT "空间 ID", MODIFY COLUMN `trace_id` varchar(255) NOT NULL COMMENT "追踪 ID", MODIFY COLUMN `node_id` varchar(255) NOT NULL COMMENT "节点 ID", MODIFY COLUMN `node_name` varchar(255) NOT NULL COMMENT "节点名称", MODIFY COLUMN `status` varchar(255) NOT NULL COMMENT "运行状态", MODIFY COLUMN `error_msg` longtext NOT NULL;
-- Create "audit" table
CREATE TABLE `audit` (
  `id` int NOT NULL,
  `user_id` int NOT NULL,
  `code` varchar(255) NOT NULL,
  `request` longtext NOT NULL,
  `response` longtext NOT NULL,
  `create_time` datetime NOT NULL,
  `url` longtext NOT NULL
) CHARSET utf8mb4 COMMENT "操作审计表";
-- Create "canvas_history" table
CREATE TABLE `canvas_history` (
  `id` int NOT NULL AUTO_INCREMENT,
  `workspace_id` varchar(255) NOT NULL,
  `draft` json NOT NULL,
  `name` varchar(255) NOT NULL,
  `create_time` datetime NOT NULL,
  `is_api` int NOT NULL COMMENT "0-不是 1-是",
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COMMENT "画布历史表" AUTO_INCREMENT 39;
-- Create "kv" table
CREATE TABLE `kv` (
  `id` int NOT NULL AUTO_INCREMENT,
  `key` varchar(255) NOT NULL,
  `value` json NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uni_key` (`key`)
) CHARSET utf8mb4 AUTO_INCREMENT 5;
-- Create "permissions" table
CREATE TABLE `permissions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(50) NOT NULL COMMENT "权限名称",
  `key` varchar(50) NOT NULL COMMENT "权限编码",
  `type` tinyint NOT NULL COMMENT "类型 1:菜单 2:按钮 3:接口",
  `parent_key` varchar(255) NULL COMMENT "父级Key",
  `path` varchar(200) NULL COMMENT "路径",
  `method` varchar(10) NULL COMMENT "HTTP方法",
  `sort` int NULL DEFAULT 0 COMMENT "排序",
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_key` (`key`)
) CHARSET utf8mb4 COMMENT "权限表" AUTO_INCREMENT 53;
-- Create "role_permissions" table
CREATE TABLE `role_permissions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `role_id` bigint unsigned NOT NULL COMMENT "角色ID",
  `permission_id` bigint unsigned NOT NULL COMMENT "权限ID",
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_role_permission` (`role_id`, `permission_id`)
) CHARSET utf8mb4 COMMENT "角色权限关联表";
-- Create "roles" table
CREATE TABLE `roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL COMMENT "角色名称",
  `code` varchar(50) NOT NULL COMMENT "角色编码",
  `description` varchar(200) NULL COMMENT "角色描述",
  `status` tinyint NOT NULL DEFAULT 1 COMMENT "状态 1:启用 0:禁用",
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_code` (`code`)
) CHARSET utf8mb4 COMMENT "角色表" AUTO_INCREMENT 33;
-- Create "user_roles" table
CREATE TABLE `user_roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL COMMENT "用户ID",
  `role_id` bigint unsigned NOT NULL COMMENT "角色ID",
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_user_role` (`user_id`, `role_id`)
) CHARSET utf8mb4 COMMENT "用户角色关联表" AUTO_INCREMENT 3;
-- Create "users" table
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL COMMENT "用户名",
  `password` varchar(100) NOT NULL COMMENT "密码",
  `salt` varchar(20) NOT NULL COMMENT "密码盐",
  `real_name` varchar(50) NULL COMMENT "真实姓名",
  `phone` varchar(20) NULL COMMENT "手机号",
  `email` varchar(100) NULL COMMENT "邮箱",
  `status` tinyint NOT NULL DEFAULT 1 COMMENT "状态 1:启用 0:禁用",
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_username` (`username`)
) CHARSET utf8mb4 COMMENT "用户表" AUTO_INCREMENT 3;
-- Drop "gogogo_kv" table
DROP TABLE `gogogo_kv`;
-- Drop "locks" table
DROP TABLE `locks`;
-- Drop "module_relation" table
DROP TABLE `module_relation`;
-- Drop "user" table
DROP TABLE `user`;

INSERT INTO permissions (`id`,`title`,`key`,`type`,`parent_key`,`path`,`method`,`sort`,`created_at`,`updated_at`) VALUES (1,'系统权限','root',1,'/',null,null,9999,'2025-01-04 17:41:19','2025-01-06 10:17:08');
INSERT INTO users (`id`,`username`,`password`,`salt`,`real_name`,`phone`,`email`,`status`,`created_at`,`updated_at`) VALUES (1,'admin','3c3d20cf4936b81600306b09ab1f6cf4','21232f297a57a5a74',null,null,null,1,'2024-12-31 13:58:23','2025-01-04 15:00:50');
INSERT INTO user_roles (`id`,`user_id`,`role_id`,`created_at`) VALUES (1,1,33,'2025-01-01 00:00:00');
INSERT INTO roles (`id`,`name`,`code`,`description`,`status`,`created_at`,`updated_at`) VALUES (33,'管理员','guanliyuan','guanliyuan',1,'2025-01-07 15:29:54','2025-01-07 15:29:54');
