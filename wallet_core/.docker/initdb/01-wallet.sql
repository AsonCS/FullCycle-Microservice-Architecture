CREATE DATABASE IF NOT EXISTS `wallet`;

USE `wallet`;

CREATE TABLE IF NOT EXISTS clients (id varchar(255), name varchar(255), email varchar(255), created_at date, updated_at date);
CREATE TABLE IF NOT EXISTS accounts (id varchar(255), client_id varchar(255), balance float, created_at date, updated_at date);
CREATE TABLE IF NOT EXISTS transactions (id varchar(255), account_id_from varchar(255), account_id_to varchar(255), amount float, created_at date);

INSERT INTO clients (id, name, email, created_at, updated_at) VALUES ('ba520d4b-e4ad-4b79-9baa-0fbaf1a1e040', 'John Doe', 'john@j.com', now(), now());
INSERT INTO clients (id, name, email, created_at, updated_at) VALUES ('12cbaad7-cea5-4b00-b689-c9fb794f01c7', 'John Doe 2', 'john2@j2.com', now(), now());

INSERT INTO accounts (id, client_id, balance, created_at, updated_at) VALUES ('f70f5ac5-0ce9-47f7-a17b-0a6b862720a9', 'ba520d4b-e4ad-4b79-9baa-0fbaf1a1e040', 1000, now(), now());
INSERT INTO accounts (id, client_id, balance, created_at, updated_at) VALUES ('8c647ed5-d414-4904-9626-d52b69dfe1bc', '12cbaad7-cea5-4b00-b689-c9fb794f01c7', 0, now(), now());