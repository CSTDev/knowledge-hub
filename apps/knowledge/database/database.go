package database

import "github.com/cstdev/knowledge-hub/apps/knowledge/types"

type Database interface {
	Create(r types.Record) error
	Search(query types.SearchQuery) ([]types.Record, error)
}
