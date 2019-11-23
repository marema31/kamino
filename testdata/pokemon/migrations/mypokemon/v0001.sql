-- +migrate Up
CREATE TABLE `pokemon` ( 
    `id` INT NOT NULL AUTO_INCREMENT , 
    `pokedex_id` INT NOT NULL ,
    `name` VARCHAR(255) NOT NULL ,
    `type1` VARCHAR(30) NOT NULL ,
    `type2` VARCHAR(30) NULL ,
    `total` INT NULL ,
    `hp` INT NULL , 
    `speed` INT NULL , 
    `generation` INT NULL , 
    `legendary` BOOLEAN NOT NULL , 
    PRIMARY KEY (`id`) 
) ENGINE = InnoDB;

-- +migrate Down
DROP TABLE `pokemon`;