CREATE TABLE `audit` (
  `id` int NOT NULL,
  `user_id` int NOT NULL,
  `code` varchar(255) NOT NULL,
  `request` longtext NOT NULL,
  `response` longtext NOT NULL,
  `create_time` datetime NOT NULL,
  `url` longtext NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作审计表';