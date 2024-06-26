// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: configs.sql

package dbc

import (
	"context"
)

const configsGet = `-- name: ConfigsGet :many
SELECT config_section
FROM public.configs
WHERE slug = ANY($1::text)
`

func (q *Queries) ConfigsGet(ctx context.Context, slugs string) ([][]byte, error) {
	rows, err := q.db.Query(ctx, configsGet, slugs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items [][]byte
	for rows.Next() {
		var config_section []byte
		if err := rows.Scan(&config_section); err != nil {
			return nil, err
		}
		items = append(items, config_section)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const saveConfig = `-- name: SaveConfig :exec
INSERT INTO public.configs (
    slug, config_section
)
VALUES (
    $1, $2
) 
ON CONFLICT (slug) DO UPDATE
SET config_section = $2
`

type SaveConfigParams struct {
	Slug          string
	ConfigSection []byte
}

func (q *Queries) SaveConfig(ctx context.Context, arg SaveConfigParams) error {
	_, err := q.db.Exec(ctx, saveConfig, arg.Slug, arg.ConfigSection)
	return err
}
