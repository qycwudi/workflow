CREATE TABLE `gogogo_kv`
(
    `id`          INT          NOT NULL,
    `spider_name` VARCHAR(255) NOT NULL,
    `k`           VARCHAR(255) NOT NULL,
    `v`           VARCHAR(255) NOT NULL,
    `timestamp`   BIGINT       NOT NULL,
    PRIMARY KEY (`id`)
);