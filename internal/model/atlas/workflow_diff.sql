-- Create "users" table
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL COMMENT "用户名",
  `password` varchar(100) NOT NULL COMMENT "密码",
  `salt` varchar(20) NOT NULL COMMENT "密码盐",
  `real_name` varchar(50) NULL COMMENT "真实姓名",
  `phone` varchar(20) NULL COMMENT "手机号",
  `email` varchar(100) NULL COMMENT "邮箱",
  `status` tinyint NOT NULL DEFAULT 1 COMMENT "状态 1:启用 0:禁用",
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_username` (`username`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT "用户表";
-- Drop "user" table
DROP TABLE `user`;
