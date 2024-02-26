package store

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/takutakahashi/billcap-schema/pkg/schema"
)

type BigQueryStore struct {
	client *bigquery.Client
	cfg    BigQueryStoreConfig
}

type BigQueryStoreConfig struct {
	Client      *bigquery.Client
	ProjectID   string
	Raw         BigQueryDatabase
	Transformed BigQueryDatabase
}

type BigQueryDatabase struct {
	DatasetID string
	TableID   string
}

func NewBigQueryStore(ctx context.Context, cfg BigQueryStoreConfig) (*BigQueryStore, error) {
	if cfg.Client != nil {
		return &BigQueryStore{client: cfg.Client, cfg: cfg}, nil
	}
	client, err := bigquery.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		return nil, err
	}
	return &BigQueryStore{client: client, cfg: cfg}, nil
}

func (s *BigQueryStore) Load(ctx context.Context, data []schema.RawData) error {
	if s.cfg.Raw.DatasetID == "" || s.cfg.Raw.TableID == "" {
		return fmt.Errorf("datasetID and tableID must be set")
	}
	inserter := s.client.Dataset(s.cfg.Raw.DatasetID).Table(s.cfg.Raw.DatasetID).Inserter()
	if err := inserter.Put(ctx, data); err != nil {
		return err
	}
	return nil
}

func (s *BigQueryStore) LoadTransformed(ctx context.Context, data []schema.TransformedData) error {
	if s.cfg.Transformed.DatasetID == "" || s.cfg.Transformed.TableID == "" {
		return fmt.Errorf("datasetID and tableID must be set")
	}
	inserter := s.client.Dataset(s.cfg.Transformed.DatasetID).Table(s.cfg.Transformed.DatasetID).Inserter()
	if err := inserter.Put(ctx, data); err != nil {
		return err
	}
	return nil
}

func (s *BigQueryStore) Transform(ctx context.Context, query string) ([]schema.TransformedData, error) {
	job, err := s.client.Query(query).Read(ctx)
	if err != nil {
		return nil, err
	}
	var transformedData []schema.TransformedData
	for {
		var row schema.TransformedData
		if err := job.Next(&row); err != nil {
			break
		}
		transformedData = append(transformedData, row)
	}
	return transformedData, nil
}

func (s *BigQueryStore) Migrate(ctx context.Context) error {
	panic("not implemented")
}

func (s *BigQueryStore) Backup(ctx context.Context) error {
	panic("not implemented")
}
