package store

import (
	"context"

	"github.com/takutakahashi/billcap-schema/pkg/schema"
)

type Store interface {
	Load(ctx context.Context, data schema.RawData) error
	Transform(ctx context.Context) ([]schema.TransformedData, error)
	Migrate(ctx context.Context) error
	Backup(ctx context.Context) error
}

type NullStore struct{}

func (s *NullStore) Load(ctx context.Context, data schema.RawData) error {
	return nil
}

func (s *NullStore) Transform(ctx context.Context) ([]schema.TransformedData, error) {
	return nil, nil
}

func (s *NullStore) Migrate(ctx context.Context) error {
	return nil
}

func (s *NullStore) Backup(ctx context.Context) error {
	return nil
}
