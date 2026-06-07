package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	es "github.com/elastic/go-elasticsearch/v8"
)

const AssetIndexName = "gm_asset"

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=client.go -destination=mock/client_mock.go -package=mock

type SearchHit struct {
	ID string
}

type SearchResponse struct {
	Total int64
	Hits  []SearchHit
}

type Client interface {
	IndexDocument(ctx context.Context, index, documentID string, body []byte) error
	Search(ctx context.Context, index string, body []byte) (*SearchResponse, error)
}

type client struct {
	es *es.Client
}

func NewClient(addresses []string) (Client, error) {
	addrs := make([]string, 0, len(addresses))
	for _, a := range addresses {
		a = strings.TrimSpace(a)
		if a != "" {
			addrs = append(addrs, a)
		}
	}
	if len(addrs) == 0 {
		return nil, fmt.Errorf("elasticsearch: no addresses configured")
	}
	c, err := es.NewClient(es.Config{Addresses: addrs})
	if err != nil {
		return nil, fmt.Errorf("elasticsearch: new client: %w", err)
	}
	return &client{es: c}, nil
}

func (c *client) IndexDocument(ctx context.Context, index, documentID string, body []byte) error {
	res, err := c.es.Index(
		index,
		strings.NewReader(string(body)),
		c.es.Index.WithContext(ctx),
		c.es.Index.WithDocumentID(documentID),
		c.es.Index.WithRefresh("false"),
	)
	if err != nil {
		return fmt.Errorf("elasticsearch index: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("elasticsearch index: %s", readErrorBody(res.Body))
	}
	return nil
}

func (c *client) Search(ctx context.Context, index string, body []byte) (*SearchResponse, error) {
	res, err := c.es.Search(
		c.es.Search.WithContext(ctx),
		c.es.Search.WithIndex(index),
		c.es.Search.WithBody(strings.NewReader(string(body))),
	)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch search: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch search: %s", readErrorBody(res.Body))
	}
	var parsed struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				ID string `json:"_id"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("elasticsearch search decode: %w", err)
	}
	out := &SearchResponse{Total: parsed.Hits.Total.Value}
	for _, h := range parsed.Hits.Hits {
		out.Hits = append(out.Hits, SearchHit{ID: h.ID})
	}
	return out, nil
}

func readErrorBody(r io.Reader) string {
	b, err := io.ReadAll(r)
	if err != nil {
		return err.Error()
	}
	return strings.TrimSpace(string(b))
}
