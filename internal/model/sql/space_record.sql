CREATE TABLE `space_record`
(
    `id`            int(11) NOT NULL AUTO_INCREMENT,
    `workspace_id`  varchar(255) NOT NULL,
    `status`        varchar(255) NOT NULL,
    `serial_number` varchar(255) NOT NULL,
    `run_time`      datetime     NOT NULL,
    `record_name`   varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `unidx_serial_number` (`serial_number`),
    KEY             `unidx_workspace_id` (`workspace_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8