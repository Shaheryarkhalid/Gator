-- +goose UP
alter TABLE feeds add COLUMN last_fetched_at timestamp ;
