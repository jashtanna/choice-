CREATE DATABASE my_new_database ;

USE my_new_database;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    company_name VARCHAR(255),
    address VARCHAR(255),
    city VARCHAR(255),
    county VARCHAR(255),
    postal VARCHAR(20),
    phone VARCHAR(20),
    email VARCHAR(255),
    web VARCHAR(255)
);
