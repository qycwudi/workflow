CREATE TABLE `locks` (
  `lock_name` varchar(255) NOT NULL COMMENT '锁名称',
  `is_locked` tinyint(1) NOT NULL COMMENT '锁是否被持有',
  `held_by` varchar(255) NOT NULL COMMENT '锁持有者',
  `locked_time` datetime NOT NULL COMMENT '锁开始持有时间',
  `timeout` int NOT NULL COMMENT '锁超时时间（秒）',
  `updated_time` datetime NOT NULL COMMENT '锁更新时间',
  `id` int NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;