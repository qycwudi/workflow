CREATE TABLE `canvas`
(
    `id`           int(11) NOT NULL AUTO_INCREMENT,
    `workspace_id` varchar(255) NOT NULL,
    `draft`        json         NOT NULL,
    `create_at`    datetime     NOT NULL,
    `update_at`    datetime     NOT NULL,
    `create_by`    varchar(255) NOT NULL,
    `update_by`    varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `unidx_workspace_id` (`workspace_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8