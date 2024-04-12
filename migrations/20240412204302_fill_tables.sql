-- +goose Up
-- +goose StatementBegin
INSERT INTO banners (feature_id, is_active)
VALUES
    (1, TRUE),
    (2, TRUE),
    (3, TRUE);

INSERT INTO banners_info (banner_id, contents)
VALUES
    (1, '{}'),
    (2, '{"content": 1}'),
    (3, '{"contentwarning": 2}');

INSERT INTO banners_tag (banner_id, tag_id)
VALUES
    (1, 1),
    (2, 1),
    (3, 2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE banners;
TRUNCATE banners_info;
TRUNCATE banners_tag;
-- +goose StatementEnd