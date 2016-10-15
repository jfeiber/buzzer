-- +goose Up
CREATE TABLE Restaurant (
  id serial PRIMARY KEY,
  name VARCHAR(99) NOT NULL,
  dateCreated timestamp NOT NULL
);

CREATE TABLE ActiveParty (
  id serial PRIMARY KEY,
  restaurantID int REFERENCES Restaurant(id),
  partyName VARCHAR(50) NOT NULL,
  partySize int NOT NULL,
  timeCreated timestamp NOT NULL,
  timeSeated timestamp,
  phoneAhead boolean NOT NULL,
  waitTimeExpected timestamp,
  waitTimeCalculated timestamp
);

CREATE TABLE Buzzer (
  id serial PRIMARY KEY,
  restaurantID int REFERENCES Restaurant(id),
  buzzerName VARCHAR(45) NOT NULL,
  lastHeartbeat timestamp,
  isActive boolean NOT NULL,
  activePartyID int REFERENCES ActiveParty(id)
);

CREATE TABLE HistoricalParty (
  id serial PRIMARY KEY,
  restaurantID int REFERENCES Restaurant(id),
  partyName VARCHAR(50) NOT NULL,
  partySize int NOT NULL,
  dateCreated timestamp NOT NULL,
  dateSeated timestamp NOT NULL,
  waitTimeExpected timestamp NOT NULL,
  waitTimeCalculated timestamp NOT NULL
);

CREATE TABLE "User" (
  id serial PRIMARY KEY,
  restaurantID int REFERENCES Restaurant(id),
  username VARCHAR(100) NOT NULL,
  password VARCHAR(100) NOT NULL,
  passSalt VARCHAR(50) NOT NULL,
  dateCreated timestamp NOT NULL
);

-- +goose Down
DROP TABLE Buzzer;
DROP TABLE HistoricalParty;
DROP TABLE "User";
DROP TABLE ActiveParty;
DROP TABLE Restaurant;
