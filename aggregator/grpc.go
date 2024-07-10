package main

import (
	"context"

	"github.com/mojtabafarzaneh/tolling/types"
)

type GRPCServer struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewGRPCServer(svc Aggregator) *GRPCServer {
	return &GRPCServer{
		svc: svc,
	}
}

func (s *GRPCServer) Aggregate(ctx context.Context, req *types.AggregateRequest) (*types.None, error) {
	distance := types.Distance{
		OBUID: int(req.ObuID),
		Value: req.Value,
		Unix:  req.Unix,
	}

	return &types.None{}, s.svc.DistanceAggregator(distance)
}
