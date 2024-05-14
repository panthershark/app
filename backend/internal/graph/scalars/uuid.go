package scalars

import (
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
)

type UUID uuid.UUID

// MarshalID: marshalls uuid
func MarshalID(id uuid.UUID) graphql.Marshaler {
	return graphql.MarshalString(id.String())
}

// UnmarshalID: converts to uuid
func UnmarshalID(v interface{}) (uuid.UUID, error) {
	switch s := v.(type) {
	case string:
		return uuid.Parse(s)
	default:
		return uuid.Nil, fmt.Errorf("%T is not a string", v)
	}
}
