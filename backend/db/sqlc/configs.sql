-- name: ConfigsGet :many
SELECT config_section
FROM public.configs
WHERE slug = ANY(sqlc.arg(slugs)::text);

-- name: SaveConfig :exec
INSERT INTO public.configs (
    slug, config_section
)
VALUES (
    sqlc.arg(slug), sqlc.arg(config_section)
) 
ON CONFLICT (slug) DO UPDATE
SET config_section = sqlc.arg(config_section);