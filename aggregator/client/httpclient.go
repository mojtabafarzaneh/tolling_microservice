package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mojtabafarzaneh/tolling/types"
)

type HTTPClient struct {
	Endpoint string
}

func NewHTTPClient(endpoint string) *HTTPClient {
	return &HTTPClient{
		Endpoint: endpoint,
	}
}

func (c *HTTPClient) Aggregate(ctx context.Context, aggreq *types.AggregateRequest) error {
	b, err := json.Marshal(aggreq)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("the service responded with non 200 statuscode %d\n", resp.StatusCode)
	}
	return nil
}
