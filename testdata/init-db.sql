SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET client_min_messages = warning;
SET row_security = off;

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA pg_catalog;

CREATE TABLE IF NOT EXISTS runners (
  id uuid NOT NULL DEFAULT uuid_generate_v1mc(),
  first_name text NOT NULL,
  last_name text NOT NULL,
  age integer,
  is_active boolean DEFAULT TRUE,
  country text NOT NULL,
  personal_best interval,
  season_best interval,
  CONSTRAINT runners_pk PRIMARY KEY (id)
);

INSERT INTO runners(first_name, last_name, age, country, personal_best, season_best)
VALUES
  ('Adam', 'Smith', 30, 'USA', '02:04:41', '02:13:13'),
  ('Sarah', 'Smith', 30, 'USA', '02:18:28', '02:18:28'),
  ('Max', 'Mueller', 28, 'Germany', '02:01:23', '01:43:21'),
  ('Julie', 'Petit', 23, 'France', '01:55:12', '01:34:34');

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
  id uuid NOT NULL DEFAULT uuid_generate_v1mc(),
  username text NOT NULL UNIQUE,
  user_password text NOT NULL,
  user_role text NOT NULL,
  access_token text,
  CONSTRAINT users_pk PRIMARY KEY (id)
);

CREATE INDEX user_access_token
ON users (access_token);

INSERT INTO users(username, user_password, user_role)
VALUES 
  ('admin', crypt('admin', gen_salt('bf')), 'admin'),
  ('user', crypt('user', gen_salt('bf')), 'user')
