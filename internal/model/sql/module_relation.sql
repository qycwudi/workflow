CREATE TABLE `module_relation`
(
    `id`        int(11) NOT NULL AUTO_INCREMENT,
    `module_id` varchar(255) NOT NULL,
    `goal_id`   varchar(255) NOT NULL,
    `types`     varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    KEY         `nidx_module_id` (`module_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8