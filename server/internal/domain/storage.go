package domain

import "context"

type StorageRepository interface {
	GenerateUploadURL(ctx context.Context, key string) (string, error)
	GenerateGetURL(ctx context.Context, key string) (string, error)
	VerifyObject(ctx context.Context, key string) (bool, error)
}
