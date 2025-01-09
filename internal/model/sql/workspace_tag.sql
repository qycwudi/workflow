CREATE TABLE `workspace_tag`
(
    `id`          int          NOT NULL AUTO_INCREMENT COMMENT '自增主建',
    `tag_name`    varchar(255) NOT NULL COMMENT '标签名称',
    `is_delete`   int          NOT NULL COMMENT '逻辑删除',
    `create_time` datetime     NOT NULL COMMENT '创建时间',
    `update_time` datetime     NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    KEY           `idx_tag_name` (`tag_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='标签表'