-- +goose Up
Create Table users (
    id UUID primary key,
    created_at timestamp,
    updated_at timestamp,
    name varchar(50) unique not null
);

-- +goose Down
DROP TABLE users;




-- goose postgres "postgres://postgres:postgres@localhost:5432/gator" down
