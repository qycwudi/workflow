CREATE TABLE `datasource` (
  `id` int NOT NULL,
  `type` varchar(255) NOT NULL,
  `config` json NOT NULL,
  `switch` int NOT NULL,
  `hash` varchar(255) NOT NULL,
  `status` varchar(255) NOT NULL,
  `create_time` datetime NOT NULL,
  `update_time` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;