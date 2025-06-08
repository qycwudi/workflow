CREATE TABLE `case` (
  `id` int NOT NULL,
  `uid` varchar(255) NOT NULL,
  `workspace_id` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `params` longtext NOT NULL,
  `create_at` datetime NOT NULL,
  `update_at` datetime NOT NULL,
  `create_by` varchar(255) NOT NULL,
  `update_by` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unidx_uid` (`uid`),
  KEY `nidx_workflow_id` (`workspace_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;