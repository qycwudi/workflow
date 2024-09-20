CREATE TABLE `workspace`
(
    `id`             int          NOT NULL AUTO_INCREMENT COMMENT '自增主建',
    `workspace_id`   varchar(255) NOT NULL COMMENT '主建',
    `workspace_name` varchar(255) NOT NULL COMMENT '名称',
    `workspace_desc` text COMMENT '描述',
    `workspace_type` varchar(50)           DEFAULT NULL COMMENT '类型flow|agent',
    `workspace_icon` varchar(255)          DEFAULT NULL COMMENT 'iconUrl',
    `canvas_config`  text COMMENT '前端画布配置',
    `create_time`    datetime     NOT NULL COMMENT '创建时间',
    `update_time`    datetime     NOT NULL COMMENT '修改时间',
    `is_delete`      int          NOT NULL DEFAULT '0' COMMENT '是否删除',
    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_workspace_id` (`workspace_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='工作空间表';