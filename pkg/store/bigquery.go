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
	ProjectID   string
	Raw         BigQueryDatabase
	Transformed BigQueryDatabase
}

type BigQueryDatabase struct {
	DatasetID string
	TableID   string
}

func NewBigQueryStore(ctx context.Context, cfg BigQueryStoreConfig) (*BigQueryStore, error) {
	client, err := bigquery.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		return nil, err
	}
	return &BigQueryStore{client: client, cfg: cfg}, nil
}

func (s *BigQueryStore) Load(ctx context.Context, data []schema.RawData) error {
	inserter := s.client.Dataset(s.cfg.Raw.DatasetID).Table(s.cfg.Raw.DatasetID).Inserter()
	if err := inserter.Put(ctx, data); err != nil {
		return err
	}
	return nil

}

func (s *BigQueryStore) Transform(ctx context.Context) ([]schema.TransformedData, error) {
	transformQuery := `

	`
	job, err := s.client.Query(transformQuery).Read(ctx)
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
	return nil
}

func (s *BigQueryStore) Backup(ctx context.Context) error {
	return nil
}
