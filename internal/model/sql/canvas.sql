CREATE TABLE `canvas`
(
    `id`             int          NOT NULL AUTO_INCREMENT,
    `canvas_id`      varchar(255) NOT NULL,
    `version`        int          NOT NULL COMMENT '版本',
    `create_time`    datetime     NOT NULL,
    `update_time`    datetime     NOT NULL,
    `debug_mode`     tinyint(1) NOT NULL COMMENT 'debug模式',
    `configuration`  json         NOT NULL COMMENT '配置信息 KV',
    `additionalInfo` json         NOT NULL COMMENT '扩展信息',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_canvas_id` (`canvas_id` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;