DROP DATABASE IF EXISTS employee;

CREATE DATABASE employee;

USE employee;

CREATE TABLE employee(
                         id int NOT NULL PRIMARY KEY AUTO_INCREMENT ,
                         age int,
                         name varchar(20) NOT NULL
);

INSERT INTO employee VALUES(1,22,'Ram');
INSERT INTO employee VALUES(2,21,'sai');
INSERT INTO employee VALUES(3,23,'kiran');