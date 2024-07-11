CREATE TABLE `gogogo_kv`
(
    `id`          int(11) NOT NULL AUTO_INCREMENT,
    `spider_name` varchar(255) NOT NULL,
    `k`           varchar(255) NOT NULL,
    `v`           text         NOT NULL,
    `timestamp`   bigint(20) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8