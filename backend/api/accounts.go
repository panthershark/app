package api

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/panthershark/app/backend/db/dbc"
)

// Account: an actor in the system. e.g. service account or user accoun
type Account struct {
	ID     uuid.UUID
	Email  string
	Person *Person
}

// Person: a person in the system. When an account has a person, then it is a user account
type Person struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	FirstName string
	LastName  string
}

// UserCreateInput: input type for creating a user
type UserCreateInput struct {
	Email     string
	FirstName string
	LastName  string
}

// AccountsApi: interfaces for working with accounts and users
type AccountsApi interface {
	GetUserById(id uuid.UUID) (*Account, error)
	GetUserByEmail(email string) (*Account, error)
	CreateAccount(email string) (uuid.UUID, error)
	CreateUser(input UserCreateInput) (uuid.UUID, error)
}

// AccountsApiFactory: injected into context so testing can mock the api.
type AccountsApiFactory func(querier dbc.Querier, ctx *context.Context) AccountsApi

type accountsApi struct {
	ctx     context.Context
	querier dbc.Querier
}

// NewAccountsApi: factor for an instance of AccountsApi
func NewAccountsApi(querier dbc.Querier, ctx *context.Context) AccountsApi {
	var c context.Context
	if ctx == nil {
		c = context.Background()
	} else {
		c = *ctx
	}

	return &accountsApi{
		ctx:     c,
		querier: querier,
	}
}

func (a *accountsApi) CreateAccount(email string) (uuid.UUID, error) {
	id, err := a.querier.CreateAccount(a.ctx, email)
	return id, err
}

func (a *accountsApi) GetUserById(id uuid.UUID) (*Account, error) {
	return nil, errors.New("not implemented")
}

func (a *accountsApi) GetUserByEmail(email string) (*Account, error) {
	rows, err := a.querier.GetAccountsByEmail(a.ctx, []string{email})
	if err != nil {
		return nil, err
	} else if len(rows) == 0 {
		return nil, errors.New("not found")
	}

	acc := Account{
		ID:    rows[0].ID,
		Email: rows[0].Email,
		Person: &Person{
			ID:        rows[0].PersonID,
			AccountID: rows[0].ID,
			FirstName: rows[0].FirstName.String,
			LastName:  rows[0].LastName.String,
		},
	}

	return &acc, nil
}

func (a *accountsApi) CreateUser(input UserCreateInput) (uuid.UUID, error) {
	accountID, err := a.querier.CreateAccount(a.ctx, input.Email)
	if err != nil {
		return uuid.Nil, err
	}

	if _, err := a.querier.CreatePerson(a.ctx, dbc.CreatePersonParams{
		AccountID: accountID,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}); err != nil {
		return uuid.Nil, err
	}

	return accountID, nil
}
