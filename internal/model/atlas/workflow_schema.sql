-- Create "api" table
CREATE TABLE `api` (`id` int NOT NULL AUTO_INCREMENT, `workspace_id` varchar(255) NOT NULL, `api_id` varchar(255) NOT NULL, `api_name` varchar(255) NOT NULL, `api_desc` text NOT NULL, `dsl` json NOT NULL, `status` varchar(255) NOT NULL, `create_time` datetime NOT NULL, `update_time` datetime NOT NULL, PRIMARY KEY (`id`), UNIQUE INDEX `unidx_api_id` (`api_id`)) CHARSET utf8mb3 COLLATE utf8mb3_general_ci COMMENT "api服务表" AUTO_INCREMENT 31;
-- Create "api_record" table
CREATE TABLE `api_record` (`id` int NOT NULL AUTO_INCREMENT, `status` varchar(255) NOT NULL, `trace_id` varchar(255) NOT NULL, `param` json NOT NULL COMMENT "参数", `extend` json NOT NULL COMMENT "扩展", `call_time` datetime NOT NULL, `api_id` varchar(255) NOT NULL, `api_name` varchar(255) NOT NULL, `error_msg` longtext NOT NULL, `secrety_key` varchar(255) NOT NULL, PRIMARY KEY (`id`)) CHARSET utf8mb3 COLLATE utf8mb3_general_ci COMMENT "api调用记录" AUTO_INCREMENT 29341;
-- Create "api_secret_key" table
CREATE TABLE `api_secret_key` (`id` int NOT NULL AUTO_INCREMENT, `secret_key` varchar(255) NOT NULL, `api_id` varchar(255) NOT NULL, `expiration_time` datetime NOT NULL, `status` varchar(255) NOT NULL COMMENT "ON、OFF", `is_deleted` int NOT NULL, `name` varchar(255) NOT NULL, PRIMARY KEY (`id`), INDEX `nidx_api_id` (`api_id`)) CHARSET utf8mb3 COLLATE utf8mb3_general_ci COMMENT "api 密钥表" AUTO_INCREMENT 57;
-- Create "audit" table
CREATE TABLE `audit` (`id` int NOT NULL, `user_id` int NOT NULL, `code` varchar(255) NOT NULL, `request` longtext NOT NULL, `response` longtext NOT NULL, `create_time` datetime NOT NULL, `url` longtext NOT NULL) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "操作审计表";
-- Create "canvas" table
CREATE TABLE `canvas` (`id` int NOT NULL AUTO_INCREMENT, `workspace_id` varchar(255) NOT NULL, `draft` json NOT NULL, `create_at` datetime NOT NULL, `update_at` datetime NOT NULL, `create_by` varchar(255) NOT NULL, `update_by` varchar(255) NOT NULL, PRIMARY KEY (`id`), UNIQUE INDEX `unidx_workspace_id` (`workspace_id`)) CHARSET utf8mb3 COLLATE utf8mb3_general_ci AUTO_INCREMENT 60;
-- Create "datasource" table
CREATE TABLE `datasource` (`id` int NOT NULL AUTO_INCREMENT, `type` varchar(255) NOT NULL, `config` json NOT NULL, `switch` int NOT NULL, `hash` varchar(255) NOT NULL, `status` varchar(255) NOT NULL, `create_time` datetime NOT NULL, `update_time` datetime NOT NULL, `name` varchar(255) NOT NULL, PRIMARY KEY (`id`), UNIQUE INDEX `un_name` (`name`) COMMENT "名称唯一索引") CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci AUTO_INCREMENT 45;
-- Create "module" table
CREATE TABLE `module` (`module_id` varchar(255) NOT NULL COMMENT "组件ID", `module_name` varchar(255) NOT NULL COMMENT "组件名称", `module_type` varchar(255) NOT NULL COMMENT "组件类型", `module_config` json NOT NULL COMMENT "组件配置", `module_index` int NOT NULL COMMENT "排序desc", PRIMARY KEY (`module_id`)) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "组件库";
-- Create "space_record" table
CREATE TABLE `space_record` (`id` int NOT NULL AUTO_INCREMENT, `workspace_id` varchar(255) NOT NULL COMMENT "空间 ID", `status` varchar(255) NOT NULL COMMENT "运行状态", `serial_number` varchar(255) NOT NULL COMMENT "流水号", `run_time` datetime NOT NULL COMMENT "运行开始时间", `record_name` varchar(255) NOT NULL COMMENT "运行记录名称", `duration` int NOT NULL COMMENT "耗时 ms", `other` json NOT NULL COMMENT "其他配置", PRIMARY KEY (`id`), UNIQUE INDEX `unidx_serial_number` (`serial_number`), INDEX `unidx_workspace_id` (`workspace_id`)) CHARSET utf8mb3 COLLATE utf8mb3_general_ci AUTO_INCREMENT 644;
-- Create "trace" table
CREATE TABLE `trace` (`id` int NOT NULL AUTO_INCREMENT, `workspace_id` varchar(255) NOT NULL COMMENT "空间 ID", `trace_id` varchar(255) NOT NULL COMMENT "追踪 ID", `input` json NOT NULL COMMENT "组件输入", `logic` json NOT NULL COMMENT "执行逻辑", `output` json NOT NULL COMMENT "组件输出", `step` int NOT NULL COMMENT "分步", `node_id` varchar(255) NOT NULL COMMENT "节点 ID", `node_name` varchar(255) NOT NULL COMMENT "节点名称", `status` varchar(255) NOT NULL COMMENT "运行状态", `elapsed_time` int NOT NULL COMMENT "运行耗时", `start_time` datetime NOT NULL COMMENT "执行时间", `error_msg` longtext NOT NULL, PRIMARY KEY (`id`), INDEX `unidx_trace_id` (`trace_id`), INDEX `unidx_workspace_id` (`workspace_id`)) CHARSET utf8mb4 COLLATE utf8mb4_general_ci COMMENT "组件链路追踪记录表" AUTO_INCREMENT 138030;
-- Create "workspace" table
CREATE TABLE `workspace` (`id` int NOT NULL AUTO_INCREMENT COMMENT "自增主建", `workspace_id` varchar(255) NOT NULL COMMENT "主建", `workspace_name` varchar(255) NOT NULL COMMENT "名称", `workspace_desc` text NULL COMMENT "描述", `workspace_type` varchar(50) NULL COMMENT "类型flow|agent", `workspace_icon` varchar(255) NULL COMMENT "iconUrl", `canvas_config` text NULL COMMENT "前端画布配置", `configuration` json NOT NULL COMMENT "配置信息 KV", `additionalInfo` json NOT NULL COMMENT "扩展信息", `create_time` datetime NOT NULL COMMENT "创建时间", `update_time` datetime NOT NULL COMMENT "修改时间", `is_delete` int NOT NULL DEFAULT 0 COMMENT "是否删除", PRIMARY KEY (`id`), UNIQUE INDEX `unique_workspace_id` (`workspace_id`)) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "工作空间表" AUTO_INCREMENT 60;
-- Create "workspace_tag" table
CREATE TABLE `workspace_tag` (`id` int NOT NULL AUTO_INCREMENT COMMENT "自增主建", `tag_name` varchar(255) NOT NULL COMMENT "标签名称", `is_delete` int NOT NULL COMMENT "逻辑删除", `create_time` datetime NOT NULL COMMENT "创建时间", `update_time` datetime NOT NULL COMMENT "修改时间", PRIMARY KEY (`id`), INDEX `idx_tag_name` (`tag_name`)) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "标签表" AUTO_INCREMENT 31;
-- Create "workspace_tag_mapping" table
CREATE TABLE `workspace_tag_mapping` (`id` int NOT NULL AUTO_INCREMENT COMMENT "主建", `tag_id` int NOT NULL COMMENT "标签ID", `workspace_id` varchar(255) NOT NULL COMMENT "画布空间ID", PRIMARY KEY (`id`), INDEX `idx_tag_id` (`tag_id`), INDEX `idx_worlspace_id` (`workspace_id`)) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "画布标签映射表" AUTO_INCREMENT 85;
-- Create "users" table
CREATE TABLE `users` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `username` varchar(50) NOT NULL COMMENT '用户名',
    `password` varchar(100) NOT NULL COMMENT '密码',
    `salt` varchar(20) NOT NULL COMMENT '密码盐',
    `real_name` varchar(50) DEFAULT NULL COMMENT '真实姓名',
    `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
    `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1:启用 0:禁用',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
-- Create "permissions" table
CREATE TABLE `permissions` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(50) NOT NULL COMMENT '权限名称',
    `code` varchar(50) NOT NULL COMMENT '权限编码',
    `type` tinyint NOT NULL COMMENT '类型 1:菜单 2:按钮 3:接口',
    `parent_id` bigint unsigned DEFAULT NULL COMMENT '父级ID',
    `path` varchar(200) DEFAULT NULL COMMENT '路径',
    `method` varchar(10) DEFAULT NULL COMMENT 'HTTP方法',
    `sort` int DEFAULT '0' COMMENT '排序',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限表';

-- Create "role_permissions" table
CREATE TABLE `role_permissions` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `role_id` bigint unsigned NOT NULL COMMENT '角色ID',
    `permission_id` bigint unsigned NOT NULL COMMENT '权限ID',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_role_permission` (`role_id`,`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限关联表';

-- Create "roles" table
CREATE TABLE `roles` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(50) NOT NULL COMMENT '角色名称',
    `code` varchar(50) NOT NULL COMMENT '角色编码',
    `description` varchar(200) DEFAULT NULL COMMENT '角色描述',
    `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态 1:启用 0:禁用',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- Create "user_roles" table
CREATE TABLE `user_roles` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `role_id` bigint unsigned NOT NULL COMMENT '角色ID',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_role` (`user_id`,`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';

-- Add initial user
INSERT INTO module (`module_id`,`module_name`,`module_type`,`module_config`,`module_index`) VALUES ('0def1c72-a83f-43a6-b6b1-c4a4c589d16b','执行代码','process','{\"type\": \"jsTransform\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"color\": \"blue\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}], \"outputs\": [{\"id\": \"Success\", \"type\": \"source\", \"color\": \"green\", \"label\": \"成功\", \"position\": \"right\", \"handleType\": \"Success\"}, {\"id\": \"Failure\", \"type\": \"source\", \"color\": \"red\", \"label\": \"失败\", \"position\": \"right\", \"handleType\": \"Failure\"}]}, \"fields\": [{\"id\": \"jsScript\", \"type\": \"code-input\", \"label\": \"处理脚本\", \"config\": {\"theme\": \"vs-dark\", \"height\": 200, \"options\": {\"minimap\": {\"enabled\": false}, \"fontSize\": 14, \"lineNumbers\": true}, \"language\": \"javascript\", \"defaultValue\": \"function Filter(msg, metadata, msgType) {  \\n  return { msg: msg, metadata: metadata, msgType: msgType };\\n}\"}}], \"runnable\": true}',5);
INSERT INTO module (`module_id`,`module_name`,`module_type`,`module_config`,`module_index`) VALUES ('0e36fd17-1f91-44d7-b124-346194e7f031','开始','input','{\"type\": \"start\", \"point\": {\"inputs\": [], \"outputs\": [{\"id\": \"Success\", \"type\": \"source\", \"color\": \"green\", \"label\": \"成功\", \"position\": \"right\", \"handleType\": \"Success\"}]}, \"fields\": [{\"id\": \"param\", \"type\": \"json-input\", \"label\": \"请求头\", \"config\": {\"height\": 150, \"defaultValue\": \"{\\\"root\\\": \\\"start\\\"}\"}}]}',1);
INSERT INTO module (`module_id`,`module_name`,`module_type`,`module_config`,`module_index`) VALUES ('11366510-4985-4db8-aab2-3c3500ec4f4e','条件','process','{\"type\": \"jsFilter\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"color\": \"blue\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}], \"outputs\": [{\"id\": \"True\", \"type\": \"source\", \"color\": \"green\", \"label\": \"True\", \"position\": \"right\", \"handleType\": \"True\"}, {\"id\": \"False\", \"type\": \"source\", \"color\": \"orange\", \"label\": \"False\", \"position\": \"right\", \"handleType\": \"False\"}, {\"id\": \"Failure\", \"type\": \"source\", \"color\": \"red\", \"label\": \"失败\", \"position\": \"right\", \"handleType\": \"Failure\"}]}, \"fields\": [{\"id\": \"jsScript\", \"type\": \"input\", \"label\": \"判断条件\", \"config\": {\"height\": 150, \"defaultValue\": \"msg.?==1\"}}]}',3);
INSERT INTO module (`module_id`,`module_name`,`module_type`,`module_config`,`module_index`) VALUES ('3c2a2245-480f-4ffb-871e-11b1389a27bf','结束','output','{\"type\": \"end\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"color\": \"blue\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}]}, \"fields\": []}',2);
INSERT INTO module (`module_id`,`module_name`,`module_type`,`module_config`,`module_index`) VALUES ('789a7fc1-9ba9-4805-9dbe-a16f69b1920d','聚合','output','{\"type\": \"join\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"color\": \"blue\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}], \"outputs\": [{\"id\": \"Success\", \"type\": \"source\", \"color\": \"green\", \"label\": \"成功\", \"position\": \"right\", \"handleType\": \"Success\"}, {\"id\": \"Failure\", \"type\": \"source\", \"color\": \"red\", \"label\": \"失败\", \"position\": \"right\", \"handleType\": \"Failure\"}]}, \"fields\": []}',8);
INSERT INTO module (`module_id`,`module_name`,`module_type`,`module_config`,`module_index`) VALUES ('83caa852-1041-4acb-8ddf-90aa7340e99d','数据库','output','{\"type\": \"database\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"color\": \"blue\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}], \"outputs\": [{\"id\": \"Success\", \"type\": \"source\", \"color\": \"green\", \"label\": \"成功\", \"position\": \"right\", \"handleType\": \"Success\"}, {\"id\": \"Failure\", \"type\": \"source\", \"color\": \"red\", \"label\": \"失败\", \"position\": \"right\", \"handleType\": \"Failure\"}]}, \"fields\": [{\"id\": \"datasource_type\", \"type\": \"radio\", \"label\": \"数据库类型\", \"config\": {\"options\": [{\"label\": \"MySQL\", \"value\": \"MySQL\"}, {\"label\": \"SqlServer\", \"value\": \"SqlServer\"}], \"defaultValue\": \"MySQL\"}}, {\"id\": \"datasource_id\", \"type\": \"number\", \"label\": \"数据源 ID\", \"config\": {\"height\": 150, \"defaultValue\": 0}}, {\"id\": \"datasource_sql\", \"type\": \"sql-input\", \"label\": \"SQL语句\", \"config\": {\"height\": 150, \"defaultValue\": \"select * from users where id = ${msg.id} and name = ${msg.name} limit 10;\"}}], \"runnable\": true}',13);
INSERT INTO module (`module_id`,`module_name`,`module_type`,`module_config`,`module_index`) VALUES ('89a43296-6f1f-4a17-b5c2-e9a2c3a7affd','HTTP-XML','output','{\"type\": \"http-xml\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"color\": \"blue\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}], \"outputs\": [{\"id\": \"Success\", \"type\": \"source\", \"color\": \"green\", \"label\": \"成功\", \"position\": \"right\", \"handleType\": \"Success\"}, {\"id\": \"Failure\", \"type\": \"source\", \"color\": \"red\", \"label\": \"失败\", \"position\": \"right\", \"handleType\": \"Failure\"}]}, \"fields\": [{\"id\": \"method\", \"type\": \"radio\", \"label\": \"请求类型\", \"config\": {\"options\": [{\"label\": \"POST\", \"value\": \"post\"}], \"defaultValue\": \"post\"}}, {\"id\": \"url\", \"type\": \"input\", \"label\": \"地址\", \"config\": {\"height\": 150, \"defaultValue\": \"http://\"}}, {\"id\": \"header\", \"type\": \"dy-form\", \"label\": \"请求头\", \"config\": {\"height\": 150, \"defaultValue\": [{\"label\": \"Content-Type\", \"value\": \"application/json\"}]}}, {\"id\": \"xmlParam\", \"type\": \"xml-input\", \"label\": \"请求参数\", \"config\": {\"theme\": \"vs-dark\", \"height\": 200, \"options\": {\"minimap\": {\"enabled\": false}, \"fontSize\": 14, \"lineNumbers\": true}, \"language\": \"xml\", \"defaultValue\": \"<?xml version=\\\"1.0\\\" encoding=\\\"UTF-8\\\"?><person><name>张三</name><age>25</age></person>\"}}], \"runnable\": true}',12);
INSERT INTO module (`module_id`,`module_name`,`module_type`,`module_config`,`module_index`) VALUES ('8d72d9aa-e4eb-40ba-94c1-8e269afb607b','HTTP','output','{\"type\": \"http\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"color\": \"blue\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}], \"outputs\": [{\"id\": \"Success\", \"type\": \"source\", \"color\": \"green\", \"label\": \"成功\", \"position\": \"right\", \"handleType\": \"Success\"}, {\"id\": \"Failure\", \"type\": \"source\", \"color\": \"red\", \"label\": \"失败\", \"position\": \"right\", \"handleType\": \"Failure\"}]}, \"fields\": [{\"id\": \"method\", \"type\": \"radio\", \"label\": \"请求类型\", \"config\": {\"options\": [{\"label\": \"GET\", \"value\": \"get\"}, {\"label\": \"POST\", \"value\": \"post\"}], \"defaultValue\": \"post\"}}, {\"id\": \"url\", \"type\": \"input\", \"label\": \"地址\", \"config\": {\"height\": 150, \"defaultValue\": \"http://\"}}, {\"id\": \"header\", \"type\": \"dy-form\", \"label\": \"请求头\", \"config\": {\"height\": 150, \"defaultValue\": [{\"label\": \"Content-Type\", \"value\": \"application/json\"}]}}], \"runnable\": true}',11);
INSERT INTO module (`module_id`,`module_name`,`module_type`,`module_config`,`module_index`) VALUES ('8fb6bee3-9e88-4c02-a25e-840ebf1f73b6','并发','input','{\"type\": \"fork\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"color\": \"blue\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}], \"outputs\": [{\"id\": \"Success\", \"type\": \"source\", \"color\": \"green\", \"label\": \"成功\", \"position\": \"right\", \"handleType\": \"Success\"}]}, \"fields\": []}',7);
INSERT INTO module (`module_id`,`module_name`,`module_type`,`module_config`,`module_index`) VALUES ('95bff99e-4b9c-4afd-94ca-9615351ef057','迭代','process','{\"type\": \"for\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"color\": \"blue\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}], \"outputs\": [{\"id\": \"Do\", \"type\": \"source\", \"color\": \"blue\", \"label\": \"迭代\", \"position\": \"right\", \"handleType\": \"Do\"}, {\"id\": \"Success\", \"type\": \"source\", \"color\": \"green\", \"label\": \"成功\", \"position\": \"right\", \"handleType\": \"Success\"}, {\"id\": \"Failure\", \"type\": \"source\", \"color\": \"red\", \"label\": \"失败\", \"position\": \"right\", \"handleType\": \"Failure\"}]}, \"fields\": [{\"id\": \"range\", \"type\": \"input\", \"label\": \"迭代对象\", \"config\": {\"height\": 150, \"defaultValue\": \"msg.?\"}}]}',4);
INSERT INTO module (`module_id`,`module_name`,`module_type`,`module_config`,`module_index`) VALUES ('bc6963b0-2fbc-4104-b9e0-3f22f1f8c8fd','文件服务器','output','{\"type\": \"fileServer\", \"point\": {\"inputs\": [{\"id\": \"input\", \"type\": \"target\", \"color\": \"blue\", \"label\": \"输入数据\", \"position\": \"left\", \"handleType\": \"Data\"}], \"outputs\": [{\"id\": \"Success\", \"type\": \"source\", \"color\": \"green\", \"label\": \"成功\", \"position\": \"right\", \"handleType\": \"Success\"}, {\"id\": \"Failure\", \"type\": \"source\", \"color\": \"red\", \"label\": \"失败\", \"position\": \"right\", \"handleType\": \"Failure\"}]}, \"fields\": [{\"id\": \"datasource_type\", \"type\": \"radio\", \"label\": \"文件系统类型\", \"config\": {\"options\": [{\"label\": \" SFTP\", \"value\": \"sftp\"}, {\"label\": \"FTP\", \"value\": \"ftp\"}], \"defaultValue\": \"ftp\"}}, {\"id\": \"datasource_mode\", \"type\": \"radio\", \"label\": \"操作模式\", \"config\": {\"options\": [{\"label\": \"上传\", \"value\": \"upload\"}, {\"label\": \"下载\", \"value\": \"download\"}, {\"label\": \"删除\", \"value\": \"delete\"}], \"defaultValue\": \"upload\"}}, {\"id\": \"datasource_id\", \"type\": \"number\", \"label\": \"数据源 ID\", \"config\": {\"height\": 150, \"defaultValue\": 0}}, {\"id\": \"datasource_path\", \"type\": \"input\", \"label\": \"路径\", \"config\": {\"height\": 150, \"defaultValue\": \"/data/xxx.xxx\"}}], \"runnable\": true}',14);
