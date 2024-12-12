CREATE TABLE `api` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `workspace_id` varchar(255) NOT NULL,
  `api_id` varchar(255) NOT NULL,
  `api_name` varchar(255) NOT NULL,
  `api_desc` text NOT NULL,
  `dsl` json NOT NULL,
  `status` varchar(255) NOT NULL,
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unidx_api_id` (`api_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 COMMENT='api服务表';