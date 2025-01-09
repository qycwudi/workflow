CREATE TABLE `permissions` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `title` varchar(50) NOT NULL COMMENT '权限名称',
    `key` varchar(50) NOT NULL COMMENT '权限编码',
    `type` tinyint NOT NULL COMMENT '类型 1:菜单 2:按钮 3:接口',
    `parent_key` varchar(255) NOT NULL DEFAULT '' COMMENT '父级权限编码',
    `path` varchar(200) DEFAULT NULL COMMENT '路径',
    `method` varchar(10) DEFAULT NULL COMMENT 'HTTP方法',
    `sort` int DEFAULT '0' COMMENT '排序',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_key` (`key`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '权限表';