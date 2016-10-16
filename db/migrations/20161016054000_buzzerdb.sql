-- +goose Up
CREATE TABLE Restaurants (
  id serial PRIMARY KEY,
  name VARCHAR(99) UNIQUE NOT NULL,
  date_created timestamp without time zone DEFAULT current_timestamp
);

CREATE TABLE ActiveParties (
  id serial PRIMARY KEY,
  restaurant_id int REFERENCES Restaurants(id),
  party_name VARCHAR(50) NOT NULL,
  party_size int NOT NULL,
  time_created timestamp NOT NULL,
  time_seated timestamp,
  phone_ahead boolean NOT NULL,
  wait_time_expected int,
  wait_time_calculated int
);

CREATE TABLE Buzzers (
  id serial PRIMARY KEY,
  restaurant_id int REFERENCES Restaurants(id),
  buzzer_name VARCHAR(45) NOT NULL,
  last_heartbeat timestamp,
  is_active boolean NOT NULL,
  activePartyID int REFERENCES ActiveParties(id)
);

CREATE TABLE HistoricalParties (
  id serial PRIMARY KEY,
  restaurant_id int REFERENCES Restaurants(id),
  party_name VARCHAR(50) NOT NULL,
  party_size int NOT NULL,
  date_created timestamp NOT NULL,
  date_seated timestamp NOT NULL,
  wait_time_expected timestamp NOT NULL,
  wait_time_calculated timestamp NOT NULL
);

CREATE TABLE Users (
  id serial PRIMARY KEY,
  restaurant_id int REFERENCES Restaurants(id),
  username VARCHAR(100) UNIQUE NOT NULL,
  password VARCHAR(100) NOT NULL,
  pass_salt VARCHAR(50) NOT NULL,
  date_created timestamp without time zone DEFAULT current_timestamp
);

-- +goose Down
DROP TABLE Buzzers;
DROP TABLE HistoricalParties;
DROP TABLE Users;
DROP TABLE ActiveParties;
DROP TABLE Restaurants;
