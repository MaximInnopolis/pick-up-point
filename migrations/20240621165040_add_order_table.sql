-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
                        id INT PRIMARY KEY,
                        user_id INT NOT NULL,
                        deadline TIMESTAMP NOT NULL,
                        is_returned BOOLEAN NOT NULL DEFAULT false,
                        is_at_pickup_point BOOLEAN NOT NULL DEFAULT false,
                        issued_to_user BOOLEAN NOT NULL DEFAULT false,
                        issued_at TIMESTAMP,
                        received_from_courier BOOLEAN NOT NULL DEFAULT true,
                        hash VARCHAR(255) NOT NULL,
                        packaging_type_id INT,
                        cost FLOAT NOT NULL,
                        FOREIGN KEY (packaging_type_id) REFERENCES packaging_types(id),
                        weight FLOAT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
