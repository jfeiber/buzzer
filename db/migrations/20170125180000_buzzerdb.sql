/*  migration to add party seated column   */

-- +goose Up
ALTER TABLE historical_parties
  ADD was_party_seated BOOLEAN NOT NULL DEFAULT 'FALSE';


-- +goose Down
ALTER TABLE historical_parties
  DROP COLUMN was_party_seated;
