//go:build tools
// +build tools

package tools

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/amacneil/dbmate/v2"
	_ "go.uber.org/mock/mockgen"
)
