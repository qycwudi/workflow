goctl model mysql ddl --src user.sql --dir .. -i ''

dump库表
kubectl exec -it xuetu-db-cc774ff4b-pd6pf  -- sh -c 'mysqldump -u root -proot wk' > mydb_dump.sql

mysql -u root -p workflow < dump.sql

#!/bin/bash

# 指定要处理的文件
DUMP_FILE="dump.sql"

# 查找并替换字符集和排序规则
sed -i 's/CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci/CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci/g' "$DUMP_FILE"
sed -i 's/COLLATE=utf8mb4_0900_ai_ci/COLLATE=utf8mb4_general_ci/g' "$DUMP_FILE"

echo "替换完成！"

chmod +x replace_collation.sh





CREATE TABLE `api`
(
    `id`           int(11) NOT NULL AUTO_INCREMENT,
    `workspace_id` varchar(255) NOT NULL,
    `api_id`       varchar(255) NOT NULL,
    `api_name`     varchar(255) NOT NULL,
    `api_desc`     text         NOT NULL,
    `dsl`          json         NOT NULL,
    `status`       varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `unidx_api_id` (`api_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='api服务表';

CREATE TABLE `api_record`
(
    `id`        int(11) NOT NULL AUTO_INCREMENT,
    `status`    varchar(255) NOT NULL,
    `trace_id`  varchar(255) NOT NULL,
    `param`     json         NOT NULL COMMENT '参数',
    `extend`    json         NOT NULL COMMENT '扩展',
    `call_time` datetime     NOT NULL,
    `api_id`    varchar(255) NOT NULL,
    `api_name`  varchar(255) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='api调用记录';

CREATE TABLE `api_secret_key`
(
    `id`              int(11) NOT NULL AUTO_INCREMENT,
    `secret_key`      varchar(255) NOT NULL,
    `api_id`          varchar(255) NOT NULL,
    `expiration_time` datetime     NOT NULL,
    PRIMARY KEY (`id`),
    KEY               `nidx_api_id` (`api_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='api 密钥表';

CREATE TABLE `canvas`
(
    `id`           int(11) NOT NULL AUTO_INCREMENT,
    `workspace_id` varchar(255) NOT NULL,
    `draft`        json         NOT NULL,
    `create_at`    datetime     NOT NULL,
    `update_at`    datetime     NOT NULL,
    `create_by`    varchar(255) NOT NULL,
    `update_by`    varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `unidx_workspace_id` (`workspace_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `gogogo_kv`
(
    `id`          int(11) NOT NULL AUTO_INCREMENT,
    `spider_name` varchar(255) NOT NULL,
    `k`           varchar(255) NOT NULL,
    `v`           text         NOT NULL,
    `timestamp`   bigint(20) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

CREATE TABLE `module`
(
    `module_id`     varchar(255) CHARACTER SET utf8mb4 NOT NULL COMMENT '组件ID',
    `module_name`   varchar(255) CHARACTER SET utf8mb4 NOT NULL COMMENT '组件名称',
    `module_type`   varchar(255) CHARACTER SET utf8mb4 NOT NULL COMMENT '组件类型',
    `module_config` json                               NOT NULL COMMENT '组件配置',
    `module_index`  int                                NOT NULL COMMENT '排序desc',
    PRIMARY KEY (`module_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组件库';

CREATE TABLE `module_relation`
(
    `id`        int(11) NOT NULL AUTO_INCREMENT,
    `module_id` varchar(255) NOT NULL,
    `goal_id`   varchar(255) NOT NULL,
    `types`     varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    KEY         `nidx_module_id` (`module_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `space_record`
(
    `id`            int(11) NOT NULL AUTO_INCREMENT,
    `workspace_id`  varchar(255) NOT NULL,
    `status`        varchar(255) NOT NULL,
    `serial_number` varchar(255) NOT NULL,
    `run_time`      datetime     NOT NULL,
    `record_name`   varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `unidx_serial_number` (`serial_number`),
    KEY             `unidx_workspace_id` (`workspace_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `trace`
(
    `id`           int(11) NOT NULL AUTO_INCREMENT,
    `workspace_id` varchar(255) NOT NULL COMMENT '空间 ID',
    `trace_id`     varchar(255) NOT NULL COMMENT '追踪 ID',
    `input`        json         NOT NULL COMMENT '组件输入',
    `logic`        json         NOT NULL COMMENT '执行逻辑',
    `output`       json         NOT NULL COMMENT '组件输出',
    `step`         int(11) NOT NULL COMMENT '分步',
    `node_id`      varchar(255) NOT NULL COMMENT '节点 ID',
    `node_name`    varchar(255) NOT NULL COMMENT '节点名称',
    `status`       varchar(255) NOT NULL COMMENT '运行状态',
    `elapsed_time` int(11) NOT NULL COMMENT '运行耗时',
    `start_time`   datetime     NOT NULL COMMENT '执行时间',
    PRIMARY KEY (`id`),
    KEY            `unidx_workspace_id` (`workspace_id`),
    KEY            `unidx_trace_id` (`trace_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='组件链路追踪记录表';

CREATE TABLE user
(
    id        bigint AUTO_INCREMENT,
    name      varchar(255) NULL COMMENT 'The username',
    password  varchar(255) NOT NULL DEFAULT '' COMMENT 'The user password',
    mobile    varchar(255) NOT NULL DEFAULT '' COMMENT 'The mobile phone number',
    gender    char(10)     NOT NULL DEFAULT 'male' COMMENT 'gender,male|female|unknown',
    nickname  varchar(255) NULL DEFAULT '' COMMENT 'The nickname',
    type      tinyint(1) NULL DEFAULT 0 COMMENT 'The user type, 0:normal,1:vip, for test golang keyword',
    create_at timestamp NULL,
    update_at timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE mobile_index (mobile),
    UNIQUE name_index (name),
    PRIMARY KEY (id)
) ENGINE = InnoDB COLLATE utf8mb4_general_ci COMMENT 'user table';

CREATE TABLE `workspace`
(
    `id`             int          NOT NULL AUTO_INCREMENT COMMENT '自增主建',
    `workspace_id`   varchar(255) NOT NULL COMMENT '主建',
    `workspace_name` varchar(255) NOT NULL COMMENT '名称',
    `workspace_desc` text COMMENT '描述',
    `workspace_type` varchar(50)           DEFAULT NULL COMMENT '类型flow|agent',
    `workspace_icon` varchar(255)          DEFAULT NULL COMMENT 'iconUrl',
    `canvas_config`  text COMMENT '前端画布配置',
    `configuration` json NOT NULL COMMENT '配置信息 KV',
    `additionalInfo` json NOT NULL COMMENT '扩展信息',
    `create_time`    datetime     NOT NULL COMMENT '创建时间',
    `update_time`    datetime     NOT NULL COMMENT '修改时间',
    `is_delete`      int          NOT NULL DEFAULT '0' COMMENT '是否删除',
    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_workspace_id` (`workspace_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='工作空间表';

CREATE TABLE `workspace_tag`
(
    `id`          int          NOT NULL AUTO_INCREMENT COMMENT '自增主建',
    `tag_name`    varchar(255) NOT NULL COMMENT '标签名称',
    `is_delete`   int          NOT NULL COMMENT '逻辑删除',
    `create_time` datetime     NOT NULL COMMENT '创建时间',
    `update_time` datetime     NOT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    KEY           `idx_tag_name` (`tag_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='标签表';

CREATE TABLE `workspace_tag_mapping`
(
    `id`           int          NOT NULL AUTO_INCREMENT COMMENT '主建',
    `tag_id`       int          NOT NULL COMMENT '标签ID',
    `workspace_id` varchar(255) NOT NULL COMMENT '画布空间ID',
    PRIMARY KEY (`id`),
    KEY            `idx_tag_id` (`tag_id`),
    KEY            `idx_worlspace_id` (`workspace_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='画布标签映射表';
