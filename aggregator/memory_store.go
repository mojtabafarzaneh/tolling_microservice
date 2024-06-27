package main

import "github.com/mojtabafarzaneh/tolling/types"

type MemoryStore struct {
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) insert(data types.Distance) error {
	return nil
}
