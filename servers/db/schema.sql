create table if not exists users (
    id int not null auto_increment primary key,
    email varchar(254) not null unique,
    username varchar(255) not null unique,
    password_hash varchar(128) not null,
    first_name varchar(64),
    last_name varchar(128),
    photo_url varchar(2083) not null,
    index email_index (email),
    index username_index (username)
);
CREATE UNIQUE INDEX unique_email_index ON users(email);
CREATE UNIQUE INDEX unique_username_index ON users(username);

create table if not exists users_signin (
    id int not null,
    date_time datetime not null,
    ip_addr varchar(255),
    constraint fk_users foreign key (id) references users(id)
);