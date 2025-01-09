CREATE TABLE `module`
(
    `module_id`     varchar(255) CHARACTER SET utf8mb4 NOT NULL COMMENT '组件ID',
    `module_name`   varchar(255) CHARACTER SET utf8mb4 NOT NULL COMMENT '组件名称',
    `module_type`   varchar(255) CHARACTER SET utf8mb4 NOT NULL COMMENT '组件类型',
    `module_config` json                               NOT NULL COMMENT '组件配置',
    `module_index`  int                                NOT NULL COMMENT '排序desc',
    PRIMARY KEY (`module_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组件库';