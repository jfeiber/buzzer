-- +goose Up
CREATE TABLE restaurants (
  id serial PRIMARY KEY,
  name VARCHAR(99) UNIQUE NOT NULL,
  date_created timestamp without time zone DEFAULT current_timestamp
);

CREATE TABLE buzzers (
  id serial PRIMARY KEY,
  restaurant_id int REFERENCES restaurants(id),
  buzzer_name VARCHAR(45) NOT NULL,
  last_heartbeat timestamp,
  is_active boolean NOT NULL
);

CREATE TABLE active_parties (
  id serial PRIMARY KEY,
  restaurant_id int REFERENCES restaurants(id),
  party_name VARCHAR(50) NOT NULL,
  party_size int NOT NULL,
  time_created timestamp NOT NULL,
  time_seated timestamp,
  phone_ahead boolean NOT NULL,
  wait_time_expected int,
  wait_time_calculated int,
  buzzer_id int REFERENCES buzzers(id)
);

CREATE TABLE historical_parties (
  id serial PRIMARY KEY,
  restaurant_id int REFERENCES restaurants(id),
  party_name VARCHAR(50) NOT NULL,
  party_size int NOT NULL,
  date_created timestamp NOT NULL,
  date_seated timestamp NOT NULL,
  wait_time_expected timestamp NOT NULL,
  wait_time_calculated timestamp NOT NULL
);

CREATE TABLE users (
  id serial PRIMARY KEY,
  restaurant_id int REFERENCES restaurants(id),
  username VARCHAR(100) UNIQUE NOT NULL,
  password VARCHAR(100) NOT NULL,
  pass_salt VARCHAR(50) NOT NULL,
  date_created timestamp without time zone DEFAULT current_timestamp
);

-- +goose Down
DROP TABLE historical_parties;
DROP TABLE users;
DROP TABLE active_parties;
DROP TABLE buzzers;
DROP TABLE restaurants;
