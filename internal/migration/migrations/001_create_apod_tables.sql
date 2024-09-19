-- +goose Up
-- +goose StatementBegin
CREATE TABLE apod_images (
 id SERIAL PRIMARY KEY,
 title TEXT NOT NULL,
 explanation TEXT NOT NULL,
 date DATE NOT NULL,
 copyright TEXT NOT NULL,
 local_storage_path TEXT NOT NULL,
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stores;
-- +goose StatementEnd

