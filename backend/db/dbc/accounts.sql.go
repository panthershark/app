// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: accounts.sql

package dbc

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createAccount = `-- name: CreateAccount :one
INSERT INTO public.accounts(email) VALUES ($1)
ON CONFLICT (email) DO NOTHING
RETURNING id
`

func (q *Queries) CreateAccount(ctx context.Context, email string) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createAccount, email)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const createPerson = `-- name: CreatePerson :one
INSERT INTO public.person(
    account_id, first_name, last_name
)
VALUES (
  $1, $2, $3
) 
RETURNING id
`

type CreatePersonParams struct {
	AccountID uuid.UUID
	FirstName string
	LastName  string
}

func (q *Queries) CreatePerson(ctx context.Context, arg CreatePersonParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createPerson, arg.AccountID, arg.FirstName, arg.LastName)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getAccountsByEmail = `-- name: GetAccountsByEmail :many
SELECT a.id, a.email, p.id::uuid as person_id, p.first_name, p.last_name
FROM public.accounts a
LEFT JOIN public.person p ON a.id = p.account_id
WHERE a.email = ANY($1::text[])
`

type GetAccountsByEmailRow struct {
	ID        uuid.UUID
	Email     string
	PersonID  uuid.UUID
	FirstName pgtype.Text
	LastName  pgtype.Text
}

func (q *Queries) GetAccountsByEmail(ctx context.Context, emails []string) ([]*GetAccountsByEmailRow, error) {
	rows, err := q.db.Query(ctx, getAccountsByEmail, emails)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*GetAccountsByEmailRow
	for rows.Next() {
		var i GetAccountsByEmailRow
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.PersonID,
			&i.FirstName,
			&i.LastName,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
