create table users (
  id bigserial NOT NULL primary key,
  name varchar(20) NOT NULL,
  password_hash varchar(64) NOT NULL,
  email varchar(255) NOT NULL
);