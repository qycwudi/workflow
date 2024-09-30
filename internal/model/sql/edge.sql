CREATE TABLE `edge`
(
    `id`          int          NOT NULL AUTO_INCREMENT,
    `edge_id`     varchar(255) NOT NULL COMMENT '边 ID',
    `edge_type`   varchar(255) NOT NULL COMMENT '边类型',
    `custom_data` varchar(255) NOT NULL COMMENT '自定义数据',
    `source`      varchar(255) NOT NULL COMMENT '起点',
    `target`      varchar(255) NOT NULL COMMENT '终点',
    `style`       json         NOT NULL COMMENT '样式',
    `route`       varchar(255) NOT NULL COMMENT '路由 True、False、Failure',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_edge_id` (`edge_id` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;