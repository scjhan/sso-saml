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

CREATE TABLE IF NOT EXISTS `idp_cluster_info` (
    `id` INT UNSIGNED AUTO_INCREMENT,
    `cid` VARCHAR(128) NOT NULL,
    `host` VARCHAR(40) NOT NULL,
    `description` VARCHAR(1024) NOT NULL,
    `pub_key` VARCHAR(1024) NOT NULL,
    PRIMARY KEY ( `id` ),
    UNIQUE KEY `host` (`host`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO idp_cluster_info(id,cid,host,description,pub_key) VALUES(NULL,"6281a69412a459016998d7b2107fd895","tm.com:7070","TM Inc.","not need now");
INSERT INTO idp_cluster_info(id,cid,host,description,pub_key) VALUES(NULL,"64e5bf6f579515fabc02cd513806a856","tb.com:6060","TB Inc.","not need now");
