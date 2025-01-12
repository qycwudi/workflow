CREATE TABLE `job` (
  `id` int NOT NULL AUTO_INCREMENT,
  `workspace_id` varchar(255) NOT NULL,
  `job_id` varchar(255) NOT NULL,
  `job_name` varchar(255) NOT NULL,
  `job_cron` varchar(255) NOT NULL,
  `job_desc` text NOT NULL,
  `params` json NOT NULL,
  `dsl` json NOT NULL,
  `status` varchar(255) NOT NULL,
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL,
  `history_id` int NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unidx_job_id` (`job_id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 COMMENT='job服务表'