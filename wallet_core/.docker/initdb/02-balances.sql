CREATE DATABASE IF NOT EXISTS `balances`;

USE `balances`;

CREATE TABLE IF NOT EXISTS accounts (
    balance float,
    id varchar(255),
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO accounts (
    balance,
    id,
    updated_at
) VALUES (
    1000,
    'f70f5ac5-0ce9-47f7-a17b-0a6b862720a9',
    NOW()
);
INSERT INTO accounts (
    balance,
    id,
    updated_at
) VALUES (
    0,
    '8c647ed5-d414-4904-9626-d52b69dfe1bc',
    NOW()
);
