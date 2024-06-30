package main

import "github.com/mojtabafarzaneh/tolling/types"

type MemoryStore struct {
	data map[int]float64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func (m *MemoryStore) insert(d types.Distance) error {
	m.data[d.OBUID] += d.Value
	return nil
}
