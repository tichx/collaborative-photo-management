CREATE TABLE IF NOT EXISTS User (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(320) NOT NULL UNIQUE,
    first_name VARCHAR(128) NOT NULL,
    last_name VARCHAR(128) NOT NULL,
    user_name VARCHAR(255) NOT NULL UNIQUE,
    pass_hash VARCHAR(128) NOT NULL,
    photo_url VARCHAR(512) NOT NULL,
    INDEX (email, user_name)
);

create table if not exists SignInLog (
    id INT NOT NULL PRIMARY KEY,
    date_time VARCHAR(320) NOT NULL,
    ip_addr VARCHAR(320) NOT NULL
);