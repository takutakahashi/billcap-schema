package store

import (
	"context"

	"cloud.google.com/go/bigquery"
	"github.com/goccy/bigquery-emulator/server"
	"github.com/goccy/bigquery-emulator/types"
	"google.golang.org/api/option"
)

func mock(ctx context.Context) (*bigquery.Client, error) {
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
}

func mockBigQueryStore() (*BigQueryStore,server.TestServer) {
	mockClient := mock()
	return &BigQueryStore{
		client: mockClient,
	}
}

func TestLoad(t *testing.T) {
	ctx := context.Background()
	client, err := mock(ctx)
	if err != nil {
		t.Error(err)
	}
}
