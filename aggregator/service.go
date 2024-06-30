package main

import (
	"fmt"

	"github.com/mojtabafarzaneh/tolling/types"
)

type Aggregator interface {
	DistanceAggregator(types.Distance) error
}

type Storer interface {
	insert(types.Distance) error
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) DistanceAggregator(data types.Distance) error {
	fmt.Printf("processing and inserting distanse in the storage %v\n", data)
	return i.store.insert(data)
}
