package main

import (
	"github.com/mojtabafarzaneh/tolling/types"
)

const basePrice = 3.15

type Aggregator interface {
	DistanceAggregator(types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}

type Storer interface {
	insert(types.Distance) error
	Get(int) (float64, error)
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
	return i.store.insert(data)
}

func (i *InvoiceAggregator) CalculateInvoice(ObuID int) (*types.Invoice, error) {
	dist, err := i.store.Get(ObuID)
	if err != nil {
		return nil, err
	}
	invoic := &types.Invoice{
		OBUID:         ObuID,
		TotalDistance: dist,
		TotalAmount:   basePrice * dist,
	}
	return invoic, nil
}
