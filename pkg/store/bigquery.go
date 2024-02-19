package store

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/takutakahashi/billcap-schema/pkg/schema"
)

type BigQueryStore struct {
	client *bigquery.Client
	cfg    BigQueryStoreConfig
}

type BigQueryStoreConfig struct {
	ProjectID            string
	RawDatasetID         string
	TransformedDatasetID string
}

func NewBigQueryStore(ctx context.Context, cfg BigQueryStoreConfig) (*BigQueryStore, error) {
	client, err := bigquery.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		return nil, err
	}
	return &BigQueryStore{client: client, cfg: cfg}, nil
}

func (s *BigQueryStore) Load(ctx context.Context, data schema.RawData) error {
	return nil
}

func (s *BigQueryStore) Transform(ctx context.Context) ([]schema.TransformedData, error) {
	return nil, nil
}

func (s *BigQueryStore) Migrate(ctx context.Context) error {
	return nil
}

func (s *BigQueryStore) Backup(ctx context.Context) error {
	return nil
}
