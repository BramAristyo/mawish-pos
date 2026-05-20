package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/BramAristyo/mawish-pos/server/internal/domain"
	appconfig "github.com/BramAristyo/mawish-pos/server/internal/infrastructure/config"
	"github.com/BramAristyo/mawish-pos/server/pkg/logger"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type R2Storage struct {
	client     *s3.Client
	presign    *s3.PresignClient
	bucketName string
	logger     *logger.ZapLogger
}

func NewR2Storage(appCfg *appconfig.Config, zapLogger *logger.ZapLogger) domain.StorageRepository {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(appCfg.R2.AccessKey, appCfg.R2.SecretAccessKey, "")),
		config.WithRegion("auto"),
	)

	if err != nil {
		zapLogger.Fatal(err.Error())
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", appCfg.R2.AccountID))
	})

	presignClient := s3.NewPresignClient(client)

	return &R2Storage{
		client:     client,
		presign:    presignClient,
		bucketName: appCfg.R2.BucketName,
		logger:     zapLogger,
	}
}

func (r *R2Storage) GenerateUploadURL(ctx context.Context, key string) (string, error) {
	presignResult, err := r.presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return "", err
	}

	return presignResult.URL, nil
}

func (r *R2Storage) GenerateGetURL(ctx context.Context, key string) (string, error) {
	presignResult, err := r.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return "", err
	}

	return presignResult.URL, nil
}

func (r *R2Storage) VerifyObject(ctx context.Context, key string) (bool, error) {
	_, err := r.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &r.bucketName,
		Key:    &key,
	})

	if err != nil {
		var notFound *types.NotFound
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
