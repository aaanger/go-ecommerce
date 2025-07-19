-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    user_email TEXT REFERENCES users(email) ON DELETE CASCADE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    status TEXT NOT NULL,
    total_price FLOAT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
