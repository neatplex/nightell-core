package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cockroachdb/errors"
	cfg "github.com/neatplex/nightell-core/internal/config"
	"github.com/neatplex/nightell-core/internal/logger"
	"io"
)

type S3 struct {
	config *cfg.Config
	l      *logger.Logger
	client *s3.Client
}

func (s *S3) Init() error {
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
		return errors.Wrap(err, "cannot load s3 config")
	} else {
		s.l.Debug("connection established with s3")
	}

	s.client = s3.NewFromConfig(c)
	return nil
}

func (s *S3) Get(path string) ([]byte, error) {
	r, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &s.config.S3.Bucket,
		Key:    &path,
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot download s3.GetObjectInput")
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	body, err := io.ReadAll(r.Body)

	return body, errors.Wrap(err, "cannot read s3.GetObjectInput.Body")
}

func (s *S3) Put(path string, body io.Reader) error {
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &s.config.S3.Bucket,
		Key:    &path,
		Body:   body,
	})
	if err != nil {
		return errors.Wrap(err, "cannot upload s3.PutObjectInput")
	}
	return nil
}

func New(c *cfg.Config, l *logger.Logger) *S3 {
	return &S3{config: c, l: l}
}
