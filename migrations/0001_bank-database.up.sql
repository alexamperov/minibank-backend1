CREATE TABLE IF NOT EXISTS users(
id serial,
username varchar(255) NOT NULL UNIQUE,
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

INSERT INTO users ( username, first_name, last_name,
                   user_role, password_hash, created_at, email)
VALUES ( "admin", "Alex", "Amperov",
        "admin","8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918",
        "2025-05-08 08:05:06", "alex.amperov@gmail.com");

INSERT INTO users ( username, first_name, last_name,
                   user_role, password_hash, created_at, email)
VALUES ( "admin", "Alex", "Amperov",
        "employee","2fdc0177057d3a5c6c2c0821e01f4fa8d90f9a3bb7afd82b0db526af98d68de8",
        "2025-05-08 08:05:06", "alex.amperov@gmail.com");

CREATE TABLE applies (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    employee_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sum NUMERIC(10, 2) NOT NULL,
    percent INT NOT NULL,
    return_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE deals (
    id SERIAL PRIMARY KEY,
    status VARCHAR(50) NOT NULL,
    sum NUMERIC(10, 2) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    percent INT NOT NULL,
    issued_at TIMESTAMP NOT NULL,
    return_at TIMESTAMP NOT NULL,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    employee_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE delays (
    id SERIAL PRIMARY KEY,
    status VARCHAR(50) NOT NULL,
    amount NUMERIC(18, 4) NOT NULL,
    accrual_date TIMESTAMP NOT NULL,
    deal_id INT NOT NULL REFERENCES deals(id) ON DELETE CASCADE,
    employee_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE pays (
    id SERIAL PRIMARY KEY,
    status VARCHAR(50) NOT NULL,
    amount NUMERIC(18, 4) NOT NULL,
    method VARCHAR(50) NOT NULL,
    deal_id INT NOT NULL REFERENCES deals(id) ON DELETE CASCADE,
    payment_date TIMESTAMP NOT NULL
);