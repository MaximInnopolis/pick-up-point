-- +goose Up
-- +goose StatementBegin
CREATE TABLE packaging_types (
                                 id SERIAL PRIMARY KEY,
                                 type VARCHAR(255) NOT NULL
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE packaging_types;
-- +goose StatementEnd
