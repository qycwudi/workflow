CREATE TABLE `space_record` (
  `id` int NOT NULL AUTO_INCREMENT,
  `workspace_id` varchar(255) CHARACTER SET utf8mb4 NOT NULL COMMENT '空间 ID',
  `status` varchar(255) CHARACTER SET utf8mb4 NOT NULL COMMENT '运行状态',
  `serial_number` varchar(255) CHARACTER SET utf8mb4 NOT NULL COMMENT '流水号',
  `run_time` datetime NOT NULL COMMENT '运行开始时间',
  `record_name` varchar(255) CHARACTER SET utf8mb4 NOT NULL COMMENT '运行记录名称',
  `duration` int NOT NULL COMMENT '耗时 ms',
  `other` json NOT NULL COMMENT '其他配置',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unidx_serial_number` (`serial_number`),
  KEY `unidx_workspace_id` (`workspace_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;