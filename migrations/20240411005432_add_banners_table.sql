-- +goose Up
-- +goose StatementBegin
CREATE TABLE banners
(
    id         BIGSERIAL PRIMARY KEY,
    feature_id INTEGER   NOT NULL,
    is_active  BOOLEAN   NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE banners_info
(
    banner_id  INTEGER,
    contents   jsonb,

    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE banners_tag
(
    banner_id INTEGER PRIMARY KEY,
    tag_id    INTEGER
);

CREATE TABLE credentials
(
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    admin BOOLEAN NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE banners_tag;
DROP TABLE banners_info;
DROP TABLE banners;
DROP TABLE credentials;
-- +goose StatementEnd
