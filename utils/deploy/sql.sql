CREATE DATABASE sso;

CREATE TABLE IF NOT EXISTS `idp_user_info`(
    `id` INT UNSIGNED AUTO_INCREMENT,
    `name` VARCHAR(40) NOT NULL,
    `passwd` VARCHAR(40) NOT NULL,
   PRIMARY KEY ( `id` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO idp_user_info(id,name,passwd) VALUES(NULL,"chenjunhan", "password");