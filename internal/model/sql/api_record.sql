CREATE TABLE `api_record`
(
    `id`        int(11) NOT NULL AUTO_INCREMENT,
    `status`    varchar(255) NOT NULL,
    `trace_id`  varchar(255) NOT NULL,
    `param`     json         NOT NULL COMMENT '参数',
    `extend`    json         NOT NULL COMMENT '扩展',
    `call_time` datetime     NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='api调用记录';