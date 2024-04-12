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
       banners_info.contents,
       banners.is_active,
       banners.created_at,
       banners_info.updated_at,
       tags_arr.tags,
       RANK() OVER (
           PARTITION BY banners_tag.tag_id, banners.feature_id
           ORDER BY banners_info.updated_at DESC
           ) rn
FROM banners
         JOIN banners_tag ON banners.id = banners_tag.banner_id
         JOIN banners_info ON banners_tab.banner_id = banners_info.banner_id
         JOIN (SELECT banners.id, array_agg(banners_tag.tag_id) as tags
               FROM banners
                        JOIN banners_tag ON banner.id = banners_tag.banner_id
               GROUP BY banners.id) tags_arr
              ON banners.id = tags_arr.id
WHERE banners_tag.tag_id = $1
   OR $1 IS NULL AND banners.feature_id = $2
   OR $2 IS NULL AND rn = 1
    LIMIT $3
OFFSET $4;

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
WHERE banners_tag.tag_id = $1 AND banners.feature_id = $2 AND banner.id >= $4
order by banners.id
LIMIT $3;

-- name: CheckBannerId :one
SELECT EXISTS(SELECT id FROM banners WHERE id = $1);

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
VALUES ($1, UNNEST($2::INT[]));

-- name: CreateBanner :one
INSERT INTO banners (feature_id, is_active, created_at)
VALUES ($1, $2, NOW())
    RETURNING id;

-- name: CreateBannerInfo :exec
INSERT INTO banners_info (banner_id, updated_at, contents)
VALUES ($1, NOW(), $2);