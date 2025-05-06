CREATE TABLE users(
id serial,
username varchar(255),
first_name varchar(255),
last_name varchar(255),
user_role varchar(25),
password_hash varchar(255),
property_type varchar(255),
address varchar(255),
phone varchar(255),
created_at timestamp,
email varchar(255)
);
