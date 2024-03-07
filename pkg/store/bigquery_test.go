package store

import (
	"context"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/goccy/bigquery-emulator/server"
	"github.com/goccy/bigquery-emulator/types"
	"github.com/takutakahashi/billcap-schema/pkg/schema"
	"google.golang.org/api/option"
)

func mock(ctx context.Context) (*bigquery.Client, *server.TestServer, error) {
	const (
		projectID = "test"
		datasetID = "dataset1"
	)
	bqServer, err := server.New(server.TempStorage)
	if err != nil {
		panic(err)
	}
	if err := bqServer.Load(
		server.StructSource(
			types.NewProject(
				projectID,
				types.NewDataset(
					datasetID,
				),
			),
		),
	); err != nil {
		panic(err)
	}
	if err := bqServer.SetProject(projectID); err != nil {
		panic(err)
	}
	testServer := bqServer.TestServer()

	client, err := bigquery.NewClient(
		ctx,
		projectID,
		option.WithEndpoint(testServer.URL),
		option.WithoutAuthentication(),
	)
	if err != nil {
		panic(err)
	}
	return client, testServer, nil
}

func mockBigQueryStore() (*BigQueryStore, *server.TestServer) {
	mockClient, server, err := mock(context.Background())
	if err != nil {
		panic(err)
	}
	return &BigQueryStore{
		client: mockClient,
		cfg: BigQueryStoreConfig{
			ProjectID: "test",
			Raw: BigQueryDatabase{
				DatasetID: "dataset1",
				TableID:   "raw",
			},
			Transformed: BigQueryDatabase{
				DatasetID: "dataset1",
				TableID:   "transformed",
			},
		},
	}, server
}

func TestLoad(t *testing.T) {
	ctx := context.Background()
	s, server := mockBigQueryStore()
	defer server.Close()
	if err := s.Load(ctx, []schema.RawData{}); err != nil {
		t.Error(err)
	}
}
