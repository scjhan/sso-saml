CREATE DATABASE sso;
use sso;

CREATE TABLE IF NOT EXISTS `idp_user_info`(
    `id` INT UNSIGNED AUTO_INCREMENT,
    `uid` VARCHAR(128) NOT NULL,
    `name` VARCHAR(40) NOT NULL,
    `passwd` VARCHAR(40) NOT NULL,
   PRIMARY KEY ( `id` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO idp_user_info(id,uid,name,passwd) VALUES(NULL,"65c2fcb282a516bd0b23c59ab137e1ba","chenjunhan", "password");
INSERT INTO idp_user_info(id,uid,name,passwd) VALUES(NULL,"9efc8fd294add9e2e12366a19f9cd1fd","gdut", "password");