-- Modify "" schema
ALTER DATABASE CHARSET utf8 COLLATE utf8_general_ci;
-- Modify "api" table
ALTER TABLE `api` CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `workspace_id` varchar(255) NOT NULL, MODIFY COLUMN `api_id` varchar(255) NOT NULL, MODIFY COLUMN `api_name` varchar(255) NOT NULL, MODIFY COLUMN `api_desc` text NOT NULL, MODIFY COLUMN `status` varchar(255) NOT NULL;
-- Modify "api_record" table
ALTER TABLE `api_record` AUTO_INCREMENT 40933, CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `status` varchar(255) NOT NULL, MODIFY COLUMN `trace_id` varchar(255) NOT NULL, MODIFY COLUMN `api_id` varchar(255) NOT NULL, MODIFY COLUMN `api_name` varchar(255) NOT NULL, MODIFY COLUMN `error_msg` longtext NOT NULL, MODIFY COLUMN `secrety_key` varchar(255) NOT NULL;
-- Modify "api_secret_key" table
ALTER TABLE `api_secret_key` CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `secret_key` varchar(255) NOT NULL, MODIFY COLUMN `api_id` varchar(255) NOT NULL, MODIFY COLUMN `status` varchar(255) NOT NULL COMMENT "ON、OFF", MODIFY COLUMN `name` varchar(255) NOT NULL;
-- Modify "canvas" table
ALTER TABLE `canvas` AUTO_INCREMENT 70, CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `workspace_id` varchar(255) NOT NULL, MODIFY COLUMN `create_by` varchar(255) NOT NULL, MODIFY COLUMN `update_by` varchar(255) NOT NULL;
-- Modify "datasource" table
ALTER TABLE `datasource` COLLATE utf8mb4_general_ci, MODIFY COLUMN `type` varchar(255) NOT NULL, MODIFY COLUMN `hash` varchar(255) NOT NULL, MODIFY COLUMN `status` varchar(255) NOT NULL, MODIFY COLUMN `name` varchar(255) NOT NULL, DROP INDEX `un_name`;
-- Modify "module" table
ALTER TABLE `module` COLLATE utf8mb4_general_ci, MODIFY COLUMN `module_id` varchar(255) NOT NULL COMMENT "组件ID", MODIFY COLUMN `module_name` varchar(255) NOT NULL COMMENT "组件名称", MODIFY COLUMN `module_type` varchar(255) NOT NULL COMMENT "组件类型";
-- Modify "space_record" table
ALTER TABLE `space_record` CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `workspace_id` varchar(255) NOT NULL COMMENT "空间 ID", MODIFY COLUMN `status` varchar(255) NOT NULL COMMENT "运行状态", MODIFY COLUMN `serial_number` varchar(255) NOT NULL COMMENT "流水号", MODIFY COLUMN `record_name` varchar(255) NOT NULL COMMENT "运行记录名称";
-- Modify "trace" table
ALTER TABLE `trace` AUTO_INCREMENT 512933;
-- Modify "users" table
ALTER TABLE `users` COMMENT "", COLLATE utf8mb4_general_ci, MODIFY COLUMN `id` int NOT NULL AUTO_INCREMENT, DROP COLUMN `username`, DROP COLUMN `password`, DROP COLUMN `salt`, DROP COLUMN `real_name`, DROP COLUMN `phone`, DROP COLUMN `email`, DROP COLUMN `status`, DROP COLUMN `created_at`, DROP COLUMN `updated_at`, ADD COLUMN `user_name` varchar(255) NOT NULL, ADD COLUMN `age` int NOT NULL;
-- Modify "workspace" table
ALTER TABLE `workspace` AUTO_INCREMENT 70, COLLATE utf8mb4_general_ci, MODIFY COLUMN `workspace_id` varchar(255) NOT NULL COMMENT "主建", MODIFY COLUMN `workspace_name` varchar(255) NOT NULL COMMENT "名称", MODIFY COLUMN `workspace_desc` text NULL COMMENT "描述", MODIFY COLUMN `workspace_type` varchar(50) NULL COMMENT "类型flow|agent", MODIFY COLUMN `workspace_icon` varchar(255) NULL COMMENT "iconUrl", MODIFY COLUMN `canvas_config` text NULL COMMENT "前端画布配置";
-- Modify "workspace_tag" table
ALTER TABLE `workspace_tag` COLLATE utf8mb4_general_ci, MODIFY COLUMN `tag_name` varchar(255) NOT NULL COMMENT "标签名称";
-- Modify "workspace_tag_mapping" table
ALTER TABLE `workspace_tag_mapping` COLLATE utf8mb4_general_ci, MODIFY COLUMN `workspace_id` varchar(255) NOT NULL COMMENT "画布空间ID";
-- Create "country" table
CREATE TABLE `country` (
  `name` varchar(255) NULL,
  `iso_code` varchar(255) NULL
) CHARSET utf8 COLLATE utf8_general_ci;
-- Create "DATABASECHANGELOG" table
CREATE TABLE `DATABASECHANGELOG` (
  `ID` varchar(255) NOT NULL,
  `AUTHOR` varchar(255) NOT NULL,
  `FILENAME` varchar(255) NOT NULL,
  `DATEEXECUTED` datetime NOT NULL,
  `ORDEREXECUTED` int NOT NULL,
  `EXECTYPE` varchar(10) NOT NULL,
  `MD5SUM` varchar(35) NULL,
  `DESCRIPTION` varchar(255) NULL,
  `COMMENTS` varchar(255) NULL,
  `TAG` varchar(255) NULL,
  `LIQUIBASE` varchar(20) NULL,
  `CONTEXTS` varchar(255) NULL,
  `LABELS` varchar(255) NULL,
  `DEPLOYMENT_ID` varchar(10) NULL,
  UNIQUE INDEX `UC_ID_AUTHOR_FILENAME` (`ID`, `AUTHOR`, `FILENAME`)
) CHARSET utf8 COLLATE utf8_general_ci;
-- Create "DATABASECHANGELOGLOCK" table
CREATE TABLE `DATABASECHANGELOGLOCK` (
  `ID` int NOT NULL,
  `LOCKED` bit NOT NULL,
  `LOCKGRANTED` datetime NULL,
  `LOCKEDBY` varchar(255) NULL,
  PRIMARY KEY (`ID`)
) CHARSET utf8 COLLATE utf8_general_ci;
-- Create "gogogo_kv" table
CREATE TABLE `gogogo_kv` (
  `id` int NOT NULL AUTO_INCREMENT,
  `spider_name` varchar(255) NOT NULL,
  `k` varchar(255) NOT NULL,
  `v` text NOT NULL,
  `timestamp` bigint NOT NULL,
  PRIMARY KEY (`id`)
) CHARSET utf8 COLLATE utf8_general_ci;
-- Create "locks" table
CREATE TABLE `locks` (
  `lock_name` varchar(255) NOT NULL COMMENT "锁名称",
  `is_locked` bool NOT NULL COMMENT "锁是否被持有",
  `held_by` varchar(255) NOT NULL COMMENT "锁持有者",
  `locked_time` datetime NOT NULL COMMENT "锁开始持有时间",
  `timeout` int NOT NULL COMMENT "锁超时时间（秒）",
  `updated_time` datetime NOT NULL COMMENT "锁更新时间",
  `id` int NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_general_ci AUTO_INCREMENT 6;
-- Create "module_relation" table
CREATE TABLE `module_relation` (
  `id` int NOT NULL AUTO_INCREMENT,
  `module_id` varchar(255) NOT NULL,
  `goal_id` varchar(255) NOT NULL,
  `types` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `nidx_module_id` (`module_id`)
) CHARSET utf8 COLLATE utf8_general_ci;
-- Create "user" table
CREATE TABLE `user` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NULL COMMENT "The username",
  `password` varchar(255) NOT NULL DEFAULT "" COMMENT "The user password",
  `mobile` varchar(255) NOT NULL DEFAULT "" COMMENT "The mobile phone number",
  `gender` char(10) NOT NULL DEFAULT "male" COMMENT "gender,male|female|unknown",
  `nickname` varchar(255) NULL DEFAULT "" COMMENT "The nickname",
  `type` bool NULL DEFAULT 0 COMMENT "The user type, 0:normal,1:vip, for test golang keyword",
  `create_at` timestamp NULL,
  `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `mobile_index` (`mobile`),
  UNIQUE INDEX `name_index` (`name`)
) CHARSET utf8mb4 COLLATE utf8mb4_general_ci COMMENT "user table";
-- Drop "audit" table
DROP TABLE `audit`;
-- Drop "permissions" table
DROP TABLE `permissions`;
-- Drop "role_permissions" table
DROP TABLE `role_permissions`;
-- Drop "roles" table
DROP TABLE `roles`;
-- Drop "user_roles" table
DROP TABLE `user_roles`;
