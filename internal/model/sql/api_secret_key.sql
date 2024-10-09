CREATE TABLE `api_secret_key`
(
    `id`              int(11) NOT NULL AUTO_INCREMENT,
    `secret_key`      varchar(255) NOT NULL,
    `api_id`          varchar(255) NOT NULL,
    `expiration_time` datetime     NOT NULL,
    PRIMARY KEY (`id`),
    KEY               `nidx_api_id` (`api_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='api 密钥表';