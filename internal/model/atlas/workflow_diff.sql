-- Modify "" schema
ALTER DATABASE CHARSET utf8 COLLATE utf8_general_ci;
-- Modify "api" table
ALTER TABLE `api` AUTO_INCREMENT 40, CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `workspace_id` varchar(255) NOT NULL, MODIFY COLUMN `api_id` varchar(255) NOT NULL, MODIFY COLUMN `api_name` varchar(255) NOT NULL, MODIFY COLUMN `api_desc` text NOT NULL, MODIFY COLUMN `status` varchar(255) NOT NULL;
-- Modify "api_record" table
ALTER TABLE `api_record` AUTO_INCREMENT 60005, CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `status` varchar(255) NOT NULL, MODIFY COLUMN `trace_id` varchar(255) NOT NULL, MODIFY COLUMN `api_id` varchar(255) NOT NULL, MODIFY COLUMN `api_name` varchar(255) NOT NULL, MODIFY COLUMN `error_msg` longtext NOT NULL, MODIFY COLUMN `secrety_key` varchar(255) NOT NULL;
-- Modify "api_secret_key" table
ALTER TABLE `api_secret_key` AUTO_INCREMENT 65, CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `secret_key` varchar(255) NOT NULL, MODIFY COLUMN `api_id` varchar(255) NOT NULL, MODIFY COLUMN `status` varchar(255) NOT NULL COMMENT "ON、OFF", MODIFY COLUMN `name` varchar(255) NOT NULL;
-- Modify "audit" table
ALTER TABLE `audit` CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `code` varchar(255) NOT NULL, MODIFY COLUMN `request` longtext NOT NULL, MODIFY COLUMN `response` longtext NOT NULL, MODIFY COLUMN `url` longtext NOT NULL;
-- Modify "canvas" table
ALTER TABLE `canvas` AUTO_INCREMENT 74, CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `workspace_id` varchar(255) NOT NULL, MODIFY COLUMN `create_by` varchar(255) NOT NULL, MODIFY COLUMN `update_by` varchar(255) NOT NULL;
-- Modify "canvas_history" table
ALTER TABLE `canvas_history` AUTO_INCREMENT 60, CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `workspace_id` varchar(255) NOT NULL, MODIFY COLUMN `name` varchar(255) NOT NULL;
-- Modify "datasource" table
ALTER TABLE `datasource` CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `type` varchar(255) NOT NULL, MODIFY COLUMN `hash` varchar(255) NOT NULL, MODIFY COLUMN `status` varchar(255) NOT NULL, MODIFY COLUMN `name` varchar(255) NOT NULL;
-- Modify "job" table
ALTER TABLE `job` AUTO_INCREMENT 5, CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `workspace_id` varchar(255) NOT NULL, MODIFY COLUMN `job_id` varchar(255) NOT NULL, MODIFY COLUMN `job_name` varchar(255) NOT NULL, MODIFY COLUMN `job_cron` varchar(255) NOT NULL, MODIFY COLUMN `job_desc` text NOT NULL, MODIFY COLUMN `status` varchar(255) NOT NULL;
-- Modify "job_record" table
ALTER TABLE `job_record` CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `status` varchar(255) NOT NULL, MODIFY COLUMN `trace_id` varchar(255) NOT NULL, MODIFY COLUMN `job_id` varchar(255) NOT NULL, MODIFY COLUMN `job_name` varchar(255) NOT NULL, MODIFY COLUMN `error_msg` longtext NOT NULL;
-- Modify "kv" table
ALTER TABLE `kv` CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `key` varchar(255) NOT NULL;
-- Modify "module" table
ALTER TABLE `module` CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `module_id` varchar(255) NOT NULL COMMENT "组件ID", MODIFY COLUMN `module_name` varchar(255) NOT NULL COMMENT "组件名称", MODIFY COLUMN `module_type` varchar(255) NOT NULL COMMENT "组件类型";
-- Modify "permissions" table
ALTER TABLE `permissions` CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `title` varchar(50) NOT NULL COMMENT "权限名称", MODIFY COLUMN `key` varchar(50) NOT NULL COMMENT "权限编码", MODIFY COLUMN `parent_key` varchar(255) NULL COMMENT "父级Key", MODIFY COLUMN `path` varchar(200) NULL COMMENT "路径", MODIFY COLUMN `method` varchar(10) NULL COMMENT "HTTP方法";
-- Modify "role_permissions" table
ALTER TABLE `role_permissions` AUTO_INCREMENT 409, CHARSET utf8, COLLATE utf8_general_ci;
-- Modify "roles" table
ALTER TABLE `roles` CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `name` varchar(50) NOT NULL COMMENT "角色名称", MODIFY COLUMN `code` varchar(50) NOT NULL COMMENT "角色编码", MODIFY COLUMN `description` varchar(200) NULL COMMENT "角色描述";
-- Modify "space_record" table
ALTER TABLE `space_record` AUTO_INCREMENT 739, CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `workspace_id` varchar(255) NOT NULL COMMENT "空间 ID", MODIFY COLUMN `status` varchar(255) NOT NULL COMMENT "运行状态", MODIFY COLUMN `serial_number` varchar(255) NOT NULL COMMENT "流水号", MODIFY COLUMN `record_name` varchar(255) NOT NULL COMMENT "运行记录名称";
-- Modify "trace" table
ALTER TABLE `trace` AUTO_INCREMENT 144671, CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `workspace_id` varchar(255) NOT NULL COMMENT "空间 ID", MODIFY COLUMN `trace_id` varchar(255) NOT NULL COMMENT "追踪 ID", MODIFY COLUMN `node_id` varchar(255) NOT NULL COMMENT "节点 ID", MODIFY COLUMN `node_name` varchar(255) NOT NULL COMMENT "节点名称", MODIFY COLUMN `status` varchar(255) NOT NULL COMMENT "运行状态", MODIFY COLUMN `error_msg` longtext NOT NULL;
-- Modify "user_roles" table
ALTER TABLE `user_roles` AUTO_INCREMENT 4, CHARSET utf8, COLLATE utf8_general_ci;
-- Modify "users" table
ALTER TABLE `users` AUTO_INCREMENT 4, CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `username` varchar(50) NOT NULL COMMENT "用户名", MODIFY COLUMN `password` varchar(100) NOT NULL COMMENT "密码", MODIFY COLUMN `salt` varchar(20) NOT NULL COMMENT "密码盐", MODIFY COLUMN `real_name` varchar(50) NULL COMMENT "真实姓名", MODIFY COLUMN `phone` varchar(20) NULL COMMENT "手机号", MODIFY COLUMN `email` varchar(100) NULL COMMENT "邮箱";
-- Modify "workspace" table
ALTER TABLE `workspace` AUTO_INCREMENT 76, CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `workspace_id` varchar(255) NOT NULL COMMENT "主建", MODIFY COLUMN `workspace_name` varchar(255) NOT NULL COMMENT "名称", MODIFY COLUMN `workspace_desc` text NULL COMMENT "描述", MODIFY COLUMN `workspace_type` varchar(50) NULL COMMENT "类型flow|agent", MODIFY COLUMN `workspace_icon` varchar(255) NULL COMMENT "iconUrl", MODIFY COLUMN `canvas_config` text NULL COMMENT "前端画布配置";
-- Modify "workspace_tag" table
ALTER TABLE `workspace_tag` AUTO_INCREMENT 35, CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `tag_name` varchar(255) NOT NULL COMMENT "标签名称";
-- Modify "workspace_tag_mapping" table
ALTER TABLE `workspace_tag_mapping` AUTO_INCREMENT 92, CHARSET utf8, COLLATE utf8_general_ci, MODIFY COLUMN `workspace_id` varchar(255) NOT NULL COMMENT "画布空间ID";
