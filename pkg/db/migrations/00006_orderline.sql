-- +goose Up
-- +goose StatementBegin
CREATE TABLE orderline (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id),
    product_id INT REFERENCES products(id),
    quantity INT,
    price FLOAT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orderline;
-- +goose StatementEnd
