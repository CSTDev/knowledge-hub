package database

import "github.com/cstdev/knowledge-hub/apps/knowledge/types"

type Database interface {
	Create(r types.Record) error
}
