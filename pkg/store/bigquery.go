package store

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/takutakahashi/billcap-schema/pkg/schema"
	"google.golang.org/api/option"
)

type BigQueryStore struct {
	client *bigquery.Client
	cfg    BigQueryStoreConfig
}

type BigQueryStoreConfig struct {
	Client          *bigquery.Client
	CredentialsPath string
	ProjectID       string
	Raw             BigQueryDatabase
	Transformed     BigQueryDatabase
}

type BigQueryDatabase struct {
	DatasetID string
	TableID   string
}

func (s *BigQueryStore) Setup(ctx context.Context) error {
	dataset := s.client.Dataset(s.cfg.Transformed.DatasetID)
	_, err := dataset.Metadata(ctx)
	if err != nil {
		if err := dataset.Create(ctx, &bigquery.DatasetMetadata{Location: "asia-northeast1"}); err != nil {
			return fmt.Errorf("failed to create dataset: %v", err)
		}
	}

	table := dataset.Table(s.cfg.Transformed.TableID)
	_, err = table.Metadata(ctx)
	if err != nil {
		schema, err := s.schema(schema.TransformedData{})
		if err != nil {
			return err
		}

		if err := table.Create(ctx, &bigquery.TableMetadata{Schema: schema}); err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
	}
	return nil

}

func NewBigQueryStore(ctx context.Context, cfg BigQueryStoreConfig) (*BigQueryStore, error) {
	if cfg.Client != nil {
		return &BigQueryStore{client: cfg.Client, cfg: cfg}, nil
	}
	client, err := bigquery.NewClient(ctx, cfg.ProjectID, option.WithCredentialsFile(cfg.CredentialsPath))
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
	valueSavers := make([]bigquery.ValueSaver, len(data))
	for i, d := range data {
		valueSavers[i] = transformedDataValueSaver{TransformedData: d}
	}

	inserter := s.client.Dataset(s.cfg.Transformed.DatasetID).Table(s.cfg.Transformed.TableID).Inserter()
	if err := inserter.Put(ctx, valueSavers); err != nil {
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

func (st *BigQueryStore) schema(s interface{}) (bigquery.Schema, error) {
	typ := reflect.TypeOf(s)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	var schema bigquery.Schema
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		var bqFieldType bigquery.FieldType
		switch field.Type {
		case reflect.TypeOf(time.Time{}), reflect.TypeOf(&time.Time{}):
			bqFieldType = bigquery.TimestampFieldType
		case reflect.TypeOf(""):
			bqFieldType = bigquery.StringFieldType
		case reflect.TypeOf(0.0):
			bqFieldType = bigquery.FloatFieldType
		default:
			return nil, fmt.Errorf("unsupported type: %s", field.Type)
		}

		schema = append(schema, &bigquery.FieldSchema{
			Name: jsonTag,
			Type: bqFieldType,
		})
	}

	return schema, nil
}

type transformedDataValueSaver struct {
	schema.TransformedData
}

func (t transformedDataValueSaver) Save() (map[string]bigquery.Value, string, error) {
	insertID := "" // Use the appropriate value for InsertID if needed.
	valueMap := make(map[string]bigquery.Value)

	// Reflect over the TransformedData struct to get the JSON tags as column names.
	tValue := reflect.ValueOf(t.TransformedData)
	tType := tValue.Type()
	for i := 0; i < tValue.NumField(); i++ {
		field := tType.Field(i)
		jsonTag := field.Tag.Get("json")
		// Do not include fields without a JSON tag.
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		valueMap[jsonTag] = tValue.Field(i).Interface()
	}

	return valueMap, insertID, nil
}
