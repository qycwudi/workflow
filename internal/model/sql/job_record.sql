CREATE TABLE `job_record` (
  `id` int NOT NULL AUTO_INCREMENT,
  `status` varchar(255) NOT NULL,
  `trace_id` varchar(255) NOT NULL,
  `param` json NOT NULL COMMENT '参数',
  `result` json NOT NULL COMMENT '结果',
  `exec_time` datetime NOT NULL,
  `job_id` varchar(255) NOT NULL,
  `job_name` varchar(255) NOT NULL,
  `error_msg` longtext NOT NULL,
  PRIMARY KEY (`id`),
  KEY `uni_job_id` (`job_id`)
) ENGINE=InnoDB COMMENT='job调用记录';