-- name: CreateAccount :one
INSERT INTO public.accounts(email) VALUES (sqlc.arg(email))
ON CONFLICT (email) DO NOTHING
RETURNING id;

-- name: CreatePerson :one
INSERT INTO public.person(
    account_id, first_name, last_name
)
VALUES (
  sqlc.arg(account_id), sqlc.arg(first_name), sqlc.arg(last_name)
) 
RETURNING id;

-- name: GetAccountsByEmail :many
SELECT a.id, a.email, p.id::uuid as person_id, p.first_name, p.last_name
FROM public.accounts a
LEFT JOIN public.person p ON a.id = p.account_id
WHERE a.email = ANY(sqlc.arg(emails)::text[]);
