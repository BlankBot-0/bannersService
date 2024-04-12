-- name: GetUserBanner :one
SELECT banners_info.contents
FROM banners
         JOIN banners_tag ON banners.id = banners_tag.banner_id
         JOIN banners_info ON banners_tag.banner_id = banners_info.banner_id
WHERE banners_tag.tag_id = $1
  AND banners.feature_id = $2
ORDER BY banners_info.updated_at DESC
LIMIT 1;

-- name: CheckActiveUserBanner :one
SELECT banners.is_active
FROM banners
         JOIN banners_tag ON banners.id = banners_tag.banner_id
WHERE banners_tag.tag_id = $1
  AND banners.feature_id = $2;

-- name: ListBanners :many
SELECT banners.id,
       banners.feature_id,
       bi.contents,
       banners.is_active,
       banners.created_at,
       bi.updated_at,
       array_agg(banners_tag.tag_id) as tags
FROM banners
         JOIN banners_tag ON banners.id = banners_tag.banner_id
         JOIN (select banner_id, contents, updated_at
               from banners_info
               order by updated_at desc
               limit 1) as bi ON banners_tag.banner_id = bi.banner_id
WHERE banners_tag.tag_id = $1
   OR $1 IS NULL AND banners.feature_id = $2
   OR $2 IS NULL AND tag_id = $1
group by 1, 2, 3, 4, 5, 6
LIMIT $3 OFFSET $4;

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
WHERE banners_tag.tag_id = $1 AND banners.feature_id = $2
LIMIT $3 OFFSET $4;

-- name: CheckBannerId :one
SELECT EXISTS(SELECT id FROM banners WHERE id = $1);

-- name: CheckExistsBanner :one
SELECT EXISTS(
    SELECT *
    FROM banners
    JOIN banners_tag ON banners.id = banners_tag.banner_id
    WHERE banners.feature_id = $1 AND banners_tag.tag_id = any(sqlc.arg(tag_ids)::INT[])
);

-- name: UpdateBannerFeature :exec
UPDATE banners
SET feature_id = $2
WHERE id = $1;

-- name: UpdateBannerIsActive :exec
UPDATE banners
SET is_active = $2
WHERE id = $1;

-- name: UpdateBannerContents :exec
INSERT INTO banners_info (banner_id, updated_at, contents)
VALUES ($1, NOW(), $2);

-- name: DeleteBannerTags :exec
DELETE
FROM banners_tag
WHERE banner_id = $1;

-- name: DeleteBannerInfo :exec
DELETE
FROM banners_info
WHERE banner_id = $1;

-- name: AddBannerTags :exec
INSERT INTO banners_tag (banner_id, tag_id)
VALUES (@banner_id::INT, UNNEST(@tag_ids::INT[]));

-- name: CreateBanner :one
INSERT INTO banners (feature_id, is_active, created_at)
VALUES ($1, $2, NOW())
RETURNING id;

-- name: CreateBannerInfo :exec
INSERT INTO banners_info (banner_id, updated_at, contents)
VALUES (@banner_id::INT, NOW(), $1);