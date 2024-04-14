-- +goose Up
-- +goose StatementBegin
CREATE INDEX tag_id_idx ON banners_tag (tag_id);
CREATE INDEX banner_id_idx ON banners (id);
CREATE INDEX feature_id_idx ON banners (feature_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX tag_id_idx;
DROP INDEX banner_id_idx;
DROP INDEX feature_id_idx;
-- +goose StatementEnd
