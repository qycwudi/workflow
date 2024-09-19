CREATE TABLE `workspace_tag_mapping`
(
    `id`           int          NOT NULL AUTO_INCREMENT COMMENT '主建',
    `tag_id`       int          NOT NULL COMMENT '标签ID',
    `workspace_id` varchar(255) NOT NULL COMMENT '画布空间ID',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='画布标签映射表'