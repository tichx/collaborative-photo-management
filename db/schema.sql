CREATE DATABASE IF NOT EXISTS PhotoCollab;

USE PhotoCollab;

CREATE TABLE Photos (
	PhotoID int not null auto_increment primary key,
	Title varchar(128) not null,
	Note varchar(320) not null,
	Tag varchar(64) not null,
	DateCreated datetime not null,
	LastModified datetime not null,
	Link varchar(512) not null,
	PhotoURL varchar(512) not null,
	UNIQUE(PhotoURL),
	UNIQUE(Link)
)

CREATE TABLE Users (
    UserID int not null auto_increment primary key,
    Username varchar(255) not null,
    PassHash varchar(255) not null,
    FirstName varchar(64) not null,
    LastName varchar(128) not null,
    AlbumID int not null,
    TagID int not null 
    Email varchar(320) not null,
    foreign key(TagID) references Tags(TagID),
    foreign key(AlbumID) references Albums(AlbumID)
    UNIQUE(Email),
    UNIQUE(Username)
);

CREATE TABLE Albums (
    AlbumID int not null auto_increment primary key,
    AlbumName varchar(255) not null,
    DateCreated datetime not null,
    PhotoID int not null,
    foreign key(PhotoID) references Photos(PhotoID)
);

CREATE TABLE Tags (
    TagID int not null auto_increment primary key,
    TagName varchar(128) not null,
    DateCreated datetime not null,
    PhotoID int not null,
    UserID int not null,
    foreign key(PhotoID) references Photos(PhotoID),
);


