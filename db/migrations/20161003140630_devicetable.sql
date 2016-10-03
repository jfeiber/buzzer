-- +goose Up
CREATE TABLE devices (
  id serial PRIMARY KEY,
  customer_id int NOT NULL,
  device_name VARCHAR(255) UNIQUE NOT NULL,
  last_heartbeat timestamp,
  is_active boolean NOT NULL,
  party_name VARCHAR(50),
  party_size int,
  wait_time int DEFAULT 0
);

-- +goose Down
DROP TABLE devices;
