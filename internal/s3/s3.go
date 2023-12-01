package s3

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	cfg "github.com/neatplex/nightel-core/internal/config"
	"go.uber.org/zap"
	"io"
)

type S3 struct {
	config *cfg.Config
	log    *zap.Logger
	client *s3.Client
}

func (s *S3) Connect() {
	credentialsCache := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		s.config.S3.AccessKey,
		s.config.S3.SecretKey,
		"",
	))

	c, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(credentialsCache),
		config.WithRegion(s.config.S3.Region),
	)
	if err != nil {
		s.log.Fatal("cannot connect to s3", zap.Error(err))
	} else {
		s.log.Debug("connection established with s3")
	}

	s.client = s3.NewFromConfig(c)
}

func (s *S3) Get(path string) ([]byte, error) {
	r, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &s.config.S3.Bucket,
		Key:    &path,
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("cannot download s3.GetObjectInput, err: %v", err.Error()))
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	body, err := io.ReadAll(r.Body)

	return body, err
}

func (s *S3) Put(path string, body io.Reader) error {
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &s.config.S3.Bucket,
		Key:    &path,
		Body:   body,
	})
	if err != nil {
		return errors.New(fmt.Sprintf("cannot upload s3.PutObjectInput, err: %v", err.Error()))
	}
	return nil
}

func New(c *cfg.Config, l *zap.Logger) *S3 {
	return &S3{config: c, log: l}
}
