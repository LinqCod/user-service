CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     name VARCHAR,
                                     surname VARCHAR,
                                     patronymic VARCHAR,
                                     age INTEGER,
                                     gender VARCHAR,
                                     nationality VARCHAR
);