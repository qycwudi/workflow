CREATE TABLE `workflow_dev`.`kv` (
	`id` INT   NOT NULL   AUTO_INCREMENT  ,
	`key` VARCHAR(255)   NOT NULL     ,
	`value` LONGTEXT   NOT NULL     ,
	UNIQUE INDEX `uni_key` (`key` ASC)  ,
	PRIMARY KEY  (`id`)  
);