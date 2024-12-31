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
