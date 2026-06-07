package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-mam/backend/internal/infra/config"
)

const (
	ConfigKeyS3Endpoint  = "MAM_STORAGE_S3_ENDPOINT"
	ConfigKeyS3AccessKey = "MAM_STORAGE_S3_ACCESS_KEY"
	ConfigKeyS3SecretKey = "MAM_STORAGE_S3_SECRET_KEY"
	ConfigKeyS3Bucket    = "MAM_STORAGE_S3_BUCKET"
	ConfigKeyS3Region    = "MAM_STORAGE_S3_REGION"
	ConfigKeyS3PathStyle = "MAM_STORAGE_S3_PATH_STYLE"
)

var (
	ErrS3EndpointMissing  = errors.New("storage s3: missing " + ConfigKeyS3Endpoint)
	ErrS3AccessKeyMissing = errors.New("storage s3: missing " + ConfigKeyS3AccessKey)
	ErrS3SecretKeyMissing = errors.New("storage s3: missing " + ConfigKeyS3SecretKey)
	ErrS3BucketMissing    = errors.New("storage s3: missing " + ConfigKeyS3Bucket)
)

type s3Settings struct {
	endpoint   string
	accessKey  string
	secretKey  string
	bucket     string
	region     string
	pathStyle  bool
}

// S3Storage implements Storage using S3-compatible APIs; rebuild via NewS3Storage when config changes.
type S3Storage struct {
	cfg config.Service

	mu        sync.Mutex
	client    *s3.Client
	presigner *s3.PresignClient
	settings  s3Settings
}

// NewS3Storage loads S3 settings from ConfigService. The AWS client is created lazily on first use.
func NewS3Storage(ctx context.Context, cfg config.Service) (*S3Storage, error) {
	settings, err := loadS3Settings(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &S3Storage{cfg: cfg, settings: settings}, nil
}

func loadS3Settings(ctx context.Context, cfg config.Service) (s3Settings, error) {
	endpoint, err := requireConfig(ctx, cfg, ConfigKeyS3Endpoint, ErrS3EndpointMissing)
	if err != nil {
		return s3Settings{}, err
	}
	accessKey, err := requireConfig(ctx, cfg, ConfigKeyS3AccessKey, ErrS3AccessKeyMissing)
	if err != nil {
		return s3Settings{}, err
	}
	secretKey, err := requireConfig(ctx, cfg, ConfigKeyS3SecretKey, ErrS3SecretKeyMissing)
	if err != nil {
		return s3Settings{}, err
	}
	bucket, err := requireConfig(ctx, cfg, ConfigKeyS3Bucket, ErrS3BucketMissing)
	if err != nil {
		return s3Settings{}, err
	}

	region, err := cfg.Get(ctx, ConfigKeyS3Region)
	if err != nil && !errors.Is(err, config.ErrNotFound) {
		return s3Settings{}, fmt.Errorf("storage s3: read %s: %w", ConfigKeyS3Region, err)
	}
	if region == "" {
		region = "us-east-1"
	}

	pathStyle := false
	pathStyleVal, err := cfg.Get(ctx, ConfigKeyS3PathStyle)
	if err != nil && !errors.Is(err, config.ErrNotFound) {
		return s3Settings{}, fmt.Errorf("storage s3: read %s: %w", ConfigKeyS3PathStyle, err)
	}
	if pathStyleVal != "" {
		pathStyle, err = strconv.ParseBool(strings.TrimSpace(pathStyleVal))
		if err != nil {
			return s3Settings{}, fmt.Errorf("storage s3: invalid %s %q: %w", ConfigKeyS3PathStyle, pathStyleVal, err)
		}
	}

	return s3Settings{
		endpoint:  strings.TrimSpace(endpoint),
		accessKey: strings.TrimSpace(accessKey),
		secretKey: strings.TrimSpace(secretKey),
		bucket:    strings.TrimSpace(bucket),
		region:    strings.TrimSpace(region),
		pathStyle: pathStyle,
	}, nil
}

func requireConfig(ctx context.Context, cfg config.Service, key string, missing error) (string, error) {
	val, err := cfg.Get(ctx, key)
	if err != nil {
		if errors.Is(err, config.ErrNotFound) {
			return "", missing
		}
		return "", fmt.Errorf("storage s3: read %s: %w", key, err)
	}
	val = strings.TrimSpace(val)
	if val == "" {
		return "", missing
	}
	return val, nil
}

func (s *S3Storage) ensureClient(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.client != nil {
		return nil
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(s.settings.region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(s.settings.accessKey, s.settings.secretKey, ""),
		),
	)
	if err != nil {
		return fmt.Errorf("storage s3: load aws config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(s.settings.endpoint)
		o.UsePathStyle = s.settings.pathStyle
	})
	s.client = client
	s.presigner = s3.NewPresignClient(client)
	return nil
}

func (s *S3Storage) Put(ctx context.Context, key string, body io.Reader, size int64, contentType string) error {
	if err := s.ensureClient(ctx); err != nil {
		return err
	}
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.settings.bucket),
		Key:    aws.String(key),
		Body:   body,
	}
	if size >= 0 {
		input.ContentLength = aws.Int64(size)
	}
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}
	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		return fmt.Errorf("storage s3: put object %q: %w", key, err)
	}
	return nil
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	if err := s.ensureClient(ctx); err != nil {
		return err
	}
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.settings.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("storage s3: delete object %q: %w", key, err)
	}
	return nil
}

func (s *S3Storage) PresignPut(ctx context.Context, key string, expire time.Duration) (string, error) {
	if err := s.ensureClient(ctx); err != nil {
		return "", err
	}
	out, err := s.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.settings.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expire))
	if err != nil {
		return "", fmt.Errorf("storage s3: presign put %q: %w", key, err)
	}
	return out.URL, nil
}