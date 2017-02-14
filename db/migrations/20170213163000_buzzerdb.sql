/*  migration to add waitlist notes column   */

-- +goose Up
ALTER TABLE active_parties
  ADD party_notes VARCHAR(150);


-- +goose Down
ALTER TABLE active_parties
  DROP COLUMN party_notes;
