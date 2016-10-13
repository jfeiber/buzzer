-- +goose Up
CREATE TABLE Device (
  id serial PRIMARY KEY,
  customer_id int NOT NULL,
  device_name VARCHAR(255) UNIQUE NOT NULL,
  last_heartbeat timestamp,
  is_active boolean NOT NULL,
  party_name VARCHAR(50),
  party_size int,
  wait_time int DEFAULT 0
);

CREATE TABLE Restaurant (
  restaurantID int PRIMARY KEY,
  name VARCHAR(99) NOT NULL,
  dateCreated timestamp NOT NULL
);

CREATE TABLE ActiveParty (
  activePartyID int PRIMARY KEY,
  restaurantID int REFERENCES Restaurant(restaurantID),
  partyName VARCHAR(50) NOT NULL,
  partySize int NOT NULL,
  timeCreated timestamp NOT NULL,
  timeSeated timestamp,
  phoneAhead boolean NOT NULL,
  waitTimeExpected timestamp,
  waitTimeCalculated timestamp
);

CREATE TABLE Buzzer (
  buzzerID serial PRIMARY KEY,
  restaurantID int REFERENCES Restaurant(restaurantID),
  buzzerName VARCHAR(45) NOT NULL,
  lastHeartbeat timestamp,
  isActive boolean NOT NULL,
  activePartyID int REFERENCES ActiveParty(activePartyID)
);

CREATE TABLE HistoricalParties (
  historicalPartiesID int PRIMARY KEY,
  restaurantID int REFERENCES Restaurant(restaurantID),
  partyName VARCHAR(50) NOT NULL,
  partySize int NOT NULL,
  dateCreated timestamp NOT NULL,
  dateSeated timestamp NOT NULL,
  waitTimeExpected timestamp NOT NULL,
  waitTimeCalculated timestamp NOT NULL
);

CREATE TABLE WebAppUser (
  webAppUserID int PRIMARY KEY,
  restaurantID int REFERENCES Restaurant(restaurantID),
  username VARCHAR(100) NOT NULL,
  password VARCHAR(100) NOT NULL,
  passSalt VARCHAR(50) NOT NULL,
  dateCreated timestamp NOT NULL
);


-- +goose Down
DROP TABLE Device;
DROP TABLE Buzzer;
DROP TABLE HistoricalParties;
DROP TABLE WebAppUser;
DROP TABLE ActiveParty;
DROP TABLE Restaurant;
