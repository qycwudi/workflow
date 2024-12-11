CREATE TABLE `api_secret_key` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `secret_key` varchar(255) NOT NULL,
  `api_id` varchar(255) NOT NULL,
  `expiration_time` datetime NOT NULL,
  `status` varchar(255) NOT NULL COMMENT 'ON、OFF',
  `is_deleted` int NOT NULL,
  PRIMARY KEY (`id`),
  KEY `nidx_api_id` (`api_id`)
) ENGINE=InnoDB COMMENT='api 密钥表';