CREATE TABLE `api_record`
(
    `id`        int(11) NOT NULL AUTO_INCREMENT,
    `status`    varchar(255) NOT NULL,
    `trace_id`  varchar(255) NOT NULL,
    `param`     json         NOT NULL COMMENT '参数',
    `extend`    json         NOT NULL COMMENT '扩展',
    `call_time` datetime     NOT NULL,
    `api_id`    varchar(255) NOT NULL,
    `api_name`  varchar(255) NOT NULL,
    `error_msg` longtext NOT NULL,
    `secrety_key` varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    KEY `uni_api_id` (`api_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='api调用记录';