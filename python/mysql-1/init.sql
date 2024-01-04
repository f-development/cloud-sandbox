create database if not exists test;
use test;
create table if not exists `test1` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `group` varchar(255) NOT NULL,
    `person` varchar(255) NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `content` varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `group_person_created_at` (`group`, `person`, `created_at`)
) ;
