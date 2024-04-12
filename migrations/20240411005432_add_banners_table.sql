-- +goose Up
-- +goose StatementBegin
-- +goose Up
-- +goose StatementBegin
CREATE TABLE banners (
id BIGSERIAL PRIMARY KEY,
feature_id INTEGER NOT NULL,
is_active BOOLEAN NOT NULL,

created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE banners_info (
banner_id INTEGER,
contents jsonb,

updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE banners_tag (
banner_id INTEGER,
tag_id INTEGER
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE banners_tag;
DROP TABLE banners_info;
DROP TABLE banners;
-- +goose StatementEnd
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE banners_tag;
DROP TABLE banners_info;
DROP TABLE banners;
-- +goose StatementEnd
