CREATE TABLE `trace` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `workspace_id` varchar(255) NOT NULL COMMENT '空间 ID',
    `trace_id` varchar(255) NOT NULL COMMENT '追踪 ID',
    `input` json NOT NULL COMMENT '组件输入',
    `logic` json NOT NULL COMMENT '执行逻辑',
    `output` json NOT NULL COMMENT '组件输出',
    `step` int(11) NOT NULL COMMENT '分步',
    `node_id` varchar(255) NOT NULL COMMENT '节点 ID',
    `node_name` varchar(255) NOT NULL COMMENT '节点名称',
    `status` varchar(255) NOT NULL COMMENT '运行状态',
    `elapsed_time` int(11) NOT NULL COMMENT '运行耗时',
    `start_time` datetime NOT NULL COMMENT '执行时间',
    `error_msg` longtext NOT NULL COMMENT '错误信息',
    PRIMARY KEY (`id`),
    KEY `unidx_workspace_id` (`workspace_id`),
    KEY `unidx_trace_id` (`trace_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT = '组件链路追踪记录表';