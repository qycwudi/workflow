-- Create "api" table
CREATE TABLE `api` (`id` int NOT NULL AUTO_INCREMENT, `workspace_id` varchar(255) NOT NULL, `api_id` varchar(255) NOT NULL, `api_name` varchar(255) NOT NULL, `api_desc` text NOT NULL, `dsl` json NOT NULL, `status` varchar(255) NOT NULL, `create_time` datetime NOT NULL, `update_time` datetime NOT NULL, `history_id` int NOT NULL, PRIMARY KEY (`id`), UNIQUE INDEX `unidx_api_id` (`api_id`)) CHARSET utf8mb3 COLLATE utf8mb3_general_ci COMMENT "api服务表" AUTO_INCREMENT 37;
-- Create "api_record" table
CREATE TABLE `api_record` (`id` int NOT NULL AUTO_INCREMENT, `status` varchar(255) NOT NULL, `trace_id` varchar(255) NOT NULL, `param` json NOT NULL COMMENT "参数", `extend` json NOT NULL COMMENT "扩展", `call_time` datetime NOT NULL, `api_id` varchar(255) NOT NULL, `api_name` varchar(255) NOT NULL, `error_msg` longtext NOT NULL, `secrety_key` varchar(255) NOT NULL, PRIMARY KEY (`id`)) CHARSET utf8mb3 COLLATE utf8mb3_general_ci COMMENT "api调用记录" AUTO_INCREMENT 59989;
-- Create "api_secret_key" table
CREATE TABLE `api_secret_key` (`id` int NOT NULL AUTO_INCREMENT, `secret_key` varchar(255) NOT NULL, `api_id` varchar(255) NOT NULL, `expiration_time` datetime NOT NULL, `status` varchar(255) NOT NULL COMMENT "ON、OFF", `is_deleted` int NOT NULL, `name` varchar(255) NOT NULL, PRIMARY KEY (`id`), INDEX `nidx_api_id` (`api_id`), INDEX `nidx_secret_key` (`secret_key`)) CHARSET utf8mb3 COLLATE utf8mb3_general_ci COMMENT "api 密钥表" AUTO_INCREMENT 63;
-- Create "audit" table
CREATE TABLE `audit` (`id` int NOT NULL, `user_id` int NOT NULL, `code` varchar(255) NOT NULL, `request` longtext NOT NULL, `response` longtext NOT NULL, `create_time` datetime NOT NULL, `url` longtext NOT NULL) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "操作审计表";
-- Create "canvas" table
CREATE TABLE `canvas` (`id` int NOT NULL AUTO_INCREMENT, `workspace_id` varchar(255) NOT NULL, `draft` json NOT NULL, `create_at` datetime NOT NULL, `update_at` datetime NOT NULL, `create_by` varchar(255) NOT NULL, `update_by` varchar(255) NOT NULL, PRIMARY KEY (`id`), UNIQUE INDEX `unidx_workspace_id` (`workspace_id`)) CHARSET utf8mb3 COLLATE utf8mb3_general_ci AUTO_INCREMENT 66;
-- Create "canvas_history" table
CREATE TABLE `canvas_history` (`id` int NOT NULL AUTO_INCREMENT, `workspace_id` varchar(255) NOT NULL, `draft` json NOT NULL, `name` varchar(255) NOT NULL, `create_time` datetime NOT NULL, `is_api` int NOT NULL COMMENT "0-不是 1-是", PRIMARY KEY (`id`)) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "画布历史表" AUTO_INCREMENT 39;
-- Create "datasource" table
CREATE TABLE `datasource` (`id` int NOT NULL AUTO_INCREMENT, `type` varchar(255) NOT NULL, `config` json NOT NULL, `switch` int NOT NULL, `hash` varchar(255) NOT NULL, `status` varchar(255) NOT NULL, `create_time` datetime NOT NULL, `update_time` datetime NOT NULL, `name` varchar(255) NOT NULL, PRIMARY KEY (`id`), UNIQUE INDEX `un_name` (`name`) COMMENT "名称唯一索引") CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci AUTO_INCREMENT 45;
-- Create "kv" table
CREATE TABLE `kv` (`id` int NOT NULL AUTO_INCREMENT, `key` varchar(255) NOT NULL, `value` json NOT NULL, PRIMARY KEY (`id`), UNIQUE INDEX `uni_key` (`key`)) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci AUTO_INCREMENT 5;
-- Create "module" table
CREATE TABLE `module` (`module_id` varchar(255) NOT NULL COMMENT "组件ID", `module_name` varchar(255) NOT NULL COMMENT "组件名称", `module_type` varchar(255) NOT NULL COMMENT "组件类型", `module_config` json NOT NULL COMMENT "组件配置", `module_index` int NOT NULL COMMENT "排序desc", PRIMARY KEY (`module_id`)) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "组件库";
-- Create "permissions" table
CREATE TABLE `permissions` (`id` bigint unsigned NOT NULL AUTO_INCREMENT, `title` varchar(50) NOT NULL COMMENT "权限名称", `key` varchar(50) NOT NULL COMMENT "权限编码", `type` tinyint NOT NULL COMMENT "类型 1:菜单 2:按钮 3:接口", `parent_key` varchar(255) NULL COMMENT "父级Key", `path` varchar(200) NULL COMMENT "路径", `method` varchar(10) NULL COMMENT "HTTP方法", `sort` int NULL DEFAULT 0 COMMENT "排序", `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (`id`), UNIQUE INDEX `idx_key` (`key`)) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "权限表" AUTO_INCREMENT 53;
-- Create "role_permissions" table
CREATE TABLE `role_permissions` (`id` bigint unsigned NOT NULL AUTO_INCREMENT, `role_id` bigint unsigned NOT NULL COMMENT "角色ID", `permission_id` bigint unsigned NOT NULL COMMENT "权限ID", `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`), UNIQUE INDEX `idx_role_permission` (`role_id`, `permission_id`)) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "角色权限关联表";
-- Create "roles" table
CREATE TABLE `roles` (`id` bigint unsigned NOT NULL AUTO_INCREMENT, `name` varchar(50) NOT NULL COMMENT "角色名称", `code` varchar(50) NOT NULL COMMENT "角色编码", `description` varchar(200) NULL COMMENT "角色描述", `status` tinyint NOT NULL DEFAULT 1 COMMENT "状态 1:启用 0:禁用", `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (`id`), UNIQUE INDEX `idx_code` (`code`)) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "角色表" AUTO_INCREMENT 33;
-- Create "space_record" table
CREATE TABLE `space_record` (`id` int NOT NULL AUTO_INCREMENT, `workspace_id` varchar(255) NOT NULL COMMENT "空间 ID", `status` varchar(255) NOT NULL COMMENT "运行状态", `serial_number` varchar(255) NOT NULL COMMENT "流水号", `run_time` datetime NOT NULL COMMENT "运行开始时间", `record_name` varchar(255) NOT NULL COMMENT "运行记录名称", `duration` int NOT NULL COMMENT "耗时 ms", `other` json NOT NULL COMMENT "其他配置", PRIMARY KEY (`id`), UNIQUE INDEX `unidx_serial_number` (`serial_number`), INDEX `unidx_workspace_id` (`workspace_id`)) CHARSET utf8mb3 COLLATE utf8mb3_general_ci AUTO_INCREMENT 670;
-- Create "trace" table
CREATE TABLE `trace` (`id` int NOT NULL AUTO_INCREMENT, `workspace_id` varchar(255) NOT NULL COMMENT "空间 ID", `trace_id` varchar(255) NOT NULL COMMENT "追踪 ID", `input` json NOT NULL COMMENT "组件输入", `logic` json NOT NULL COMMENT "执行逻辑", `output` json NOT NULL COMMENT "组件输出", `step` int NOT NULL COMMENT "分步", `node_id` varchar(255) NOT NULL COMMENT "节点 ID", `node_name` varchar(255) NOT NULL COMMENT "节点名称", `status` varchar(255) NOT NULL COMMENT "运行状态", `elapsed_time` int NOT NULL COMMENT "运行耗时", `start_time` datetime NOT NULL COMMENT "执行时间", `error_msg` longtext NOT NULL, PRIMARY KEY (`id`), INDEX `unidx_trace_id` (`trace_id`), INDEX `unidx_workspace_id` (`workspace_id`)) CHARSET utf8mb4 COLLATE utf8mb4_general_ci COMMENT "组件链路追踪记录表" AUTO_INCREMENT 144482;
-- Create "user_roles" table
CREATE TABLE `user_roles` (`id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL COMMENT "用户ID", `role_id` bigint unsigned NOT NULL COMMENT "角色ID", `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`), UNIQUE INDEX `idx_user_role` (`user_id`, `role_id`)) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "用户角色关联表" AUTO_INCREMENT 3;
-- Create "users" table
CREATE TABLE `users` (`id` bigint unsigned NOT NULL AUTO_INCREMENT, `username` varchar(50) NOT NULL COMMENT "用户名", `password` varchar(100) NOT NULL COMMENT "密码", `salt` varchar(20) NOT NULL COMMENT "密码盐", `real_name` varchar(50) NULL COMMENT "真实姓名", `phone` varchar(20) NULL COMMENT "手机号", `email` varchar(100) NULL COMMENT "邮箱", `status` tinyint NOT NULL DEFAULT 1 COMMENT "状态 1:启用 0:禁用", `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (`id`), UNIQUE INDEX `idx_username` (`username`)) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "用户表" AUTO_INCREMENT 3;
-- Create "workspace" table
CREATE TABLE `workspace` (`id` int NOT NULL AUTO_INCREMENT COMMENT "自增主建", `workspace_id` varchar(255) NOT NULL COMMENT "主建", `workspace_name` varchar(255) NOT NULL COMMENT "名称", `workspace_desc` text NULL COMMENT "描述", `workspace_type` varchar(50) NULL COMMENT "类型flow|agent", `workspace_icon` varchar(255) NULL COMMENT "iconUrl", `canvas_config` text NULL COMMENT "前端画布配置", `configuration` json NOT NULL COMMENT "配置信息 KV", `additionalInfo` json NOT NULL COMMENT "扩展信息", `create_time` datetime NOT NULL COMMENT "创建时间", `update_time` datetime NOT NULL COMMENT "修改时间", `is_delete` int NOT NULL DEFAULT 0 COMMENT "是否删除", PRIMARY KEY (`id`), UNIQUE INDEX `unique_workspace_id` (`workspace_id`)) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "工作空间表" AUTO_INCREMENT 68;
-- Create "workspace_tag" table
CREATE TABLE `workspace_tag` (`id` int NOT NULL AUTO_INCREMENT COMMENT "自增主建", `tag_name` varchar(255) NOT NULL COMMENT "标签名称", `is_delete` int NOT NULL COMMENT "逻辑删除", `create_time` datetime NOT NULL COMMENT "创建时间", `update_time` datetime NOT NULL COMMENT "修改时间", PRIMARY KEY (`id`), INDEX `idx_tag_name` (`tag_name`)) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "标签表" AUTO_INCREMENT 33;
-- Create "workspace_tag_mapping" table
CREATE TABLE `workspace_tag_mapping` (`id` int NOT NULL AUTO_INCREMENT COMMENT "主建", `tag_id` int NOT NULL COMMENT "标签ID", `workspace_id` varchar(255) NOT NULL COMMENT "画布空间ID", PRIMARY KEY (`id`), INDEX `idx_tag_id` (`tag_id`), INDEX `idx_worlspace_id` (`workspace_id`)) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "画布标签映射表" AUTO_INCREMENT 90;
