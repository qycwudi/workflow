CREATE TABLE `node`
(
    `id`                   int                                                           NOT NULL AUTO_INCREMENT,
    `node_id`              varchar(255)                                                  NOT NULL,
    `node_type`            varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '组件类型',
    `label_config`         json                                                          NOT NULL COMMENT '前端字段配置',
    `custom_config`        json                                                          NOT NULL COMMENT '组件自定义配置',
    `task_config`          json                                                          NOT NULL COMMENT '任务配置',
    `style_config`         json                                                          NOT NULL COMMENT '样式配置',
    `anchor_points_config` json                                                          NOT NULL COMMENT '锚点配置',
    `position`             json                                                          NOT NULL COMMENT '坐标配置',
    `create_time`          datetime                                                      NOT NULL,
    `update_time`          datetime                                                      NOT NULL,
    `node_name`            varchar(255)                                                  NOT NULL COMMENT '节点名称',
    `configuration`        json                                                          NOT NULL COMMENT '组件通用配置',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_node_id` (`node_id` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;