package database

import "github.com/cstdev/knowledge-hub/apps/knowledge/types"

type Database interface {
	Create(r types.Record) error
	Search(query types.SearchQuery) ([]types.Record, error)
	Update(id int, r types.Record) error
	Delete(id int) error
}

type FakeDB struct {
}

func (f *FakeDB) Create(r types.Record) error {
	return nil
}

func (f *FakeDB) Search(query types.SearchQuery) ([]types.Record, error) {
	return nil, nil
}

func (f *FakeDB) Update(id int, r types.Record) error {
	return nil
}

func (f *FakeDB) Delete(id int) error {
	return nil
}
