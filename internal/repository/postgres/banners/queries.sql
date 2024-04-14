-- name: GetUserBanner :one
SELECT banners_info.contents
FROM banners
         JOIN banners_tag ON banners.id = banners_tag.banner_id
         JOIN banners_info ON banners_tag.banner_id = banners_info.banner_id
WHERE banners_tag.tag_id = @tag_id::INT
  AND banners.feature_id = @feature_id::INT
ORDER BY banners_info.updated_at DESC
LIMIT 1;

-- name: CheckActiveUserBanner :one
SELECT banners.is_active
FROM banners
         JOIN banners_tag ON banners.id = banners_tag.banner_id
WHERE banners_tag.tag_id = @tag_id::INT
  AND banners.feature_id = @feature_id::INT;

-- name: ListBanners :many
SELECT banners.id,
       banners.feature_id,
       bi.contents,
       banners.is_active,
       banners.created_at,
       bi.updated_at,
       array_agg(banners_tag.tag_id)::INT[] as tags
FROM banners
         JOIN banners_tag ON banners.id = banners_tag.banner_id
         JOIN (SELECT banners_info.banner_id, banners_info.contents, banners_info.updated_at
               FROM banners_info
                        RIGHT JOIN (SELECT banner_id, MAX(updated_at) as upd
                                    FROM banners_info
                                    GROUP BY banner_id) as upds ON banners_info.banner_id = upds.banner_id
                   AND banners_info.updated_at = upds.upd) as bi
              ON banners_tag.banner_id = bi.banner_id
WHERE banners_tag.tag_id = sqlc.narg(tag_id)
   OR sqlc.narg(tag_id) IS NULL
    AND banners.feature_id = sqlc.narg(feature_id)
   OR sqlc.narg(feature_id) IS NULL
group by 1, 2, 3, 4, 5, 6
LIMIT @limit_val::INT OFFSET @offset_val::INT;

-- name: ListBannerVersions :many
SELECT banners.id,
       banners.feature_id,
       banners_info.contents,
       banners.is_active,
       banners.created_at,
       banners_info.updated_at
FROM banners
         JOIN banners_tag ON banners.id = banners_tag.banner_id
         JOIN banners_info ON banners_tag.banner_id = banners_info.banner_id
WHERE banners_tag.tag_id = @tag_id::INT
  AND banners.feature_id = @feature_id::INT
LIMIT @limit_val::INT OFFSET @offset_val::INT;

-- name: CheckBannerId :one
SELECT EXISTS(SELECT id FROM banners WHERE id = @banner_id::INT);

-- name: CheckExistsBanner :one
SELECT EXISTS(SELECT *
              FROM banners
                       JOIN banners_tag ON banners.id = banners_tag.banner_id
              WHERE banners.feature_id = @feature_id::INT
                AND banners_tag.tag_id = any (@tag_ids::INT[]));

-- name: UpdateBannerFeature :exec
UPDATE banners
SET feature_id = @feature_id::INT
WHERE id = @banner_id::INT;

-- name: UpdateBannerIsActive :exec
UPDATE banners
SET is_active = @is_active::BOOLEAN
WHERE id = @banner_id::INT;

-- name: UpdateBannerContents :exec
INSERT INTO banners_info (banner_id, updated_at, contents)
VALUES (@banner_id::INT, NOW(), @contents);

-- name: DeleteBannerTags :exec
DELETE
FROM banners_tag
WHERE banner_id = @banner_id::INT;

-- name: DeleteBannerInfo :exec
DELETE
FROM banners_info
WHERE banner_id = @banner_id::INT;

-- name: DeleteBanner :exec
DELETE
FROM banners
WHERE id = @banner_id::INT;

-- name: AddBannerTags :exec
INSERT INTO banners_tag (banner_id, tag_id)
VALUES (@banner_id::INT, UNNEST(@tag_ids::INT[]));

-- name: CreateBanner :one
INSERT INTO banners (feature_id, is_active, created_at)
VALUES (@feature_id::INT, @is_active::BOOLEAN, NOW())
RETURNING id;

-- name: CreateBannerInfo :exec
INSERT INTO banners_info (banner_id, updated_at, contents)
VALUES (@banner_id::INT, NOW(), @contents);