CREATE TABLE `workspace_tag_mapping`
(
    `id`           int          NOT NULL AUTO_INCREMENT COMMENT '主建',
    `tag_id`       int          NOT NULL COMMENT '标签ID',
    `workspace_id` varchar(255) NOT NULL COMMENT '画布空间ID',
    PRIMARY KEY (`id`),
    KEY            `idx_tag_id` (`tag_id`),
    KEY            `idx_worlspace_id` (`workspace_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='画布标签映射表'