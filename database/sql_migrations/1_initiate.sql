-- +migrate Up
-- +migrate StatementBegin

create table cinema (
    id          BIGINT NOT NULL,
    name  varchar(256),
    location   varchar(256),
    rating   float
)

-- +migrate StatementEnd