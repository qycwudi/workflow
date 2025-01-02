CREATE TABLE `canvas_history` (
	`id` INT   NOT NULL   AUTO_INCREMENT  ,
	`workspace_id` VARCHAR(255)   NOT NULL     ,
	`draft` JSON   NOT NULL     ,
	`name` VARCHAR(255)   NOT NULL     ,
	`create_time` DATETIME   NOT NULL     ,
	PRIMARY KEY  (`id`)  
) COMMENT='画布历史表';