package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/logging"
	"github.com/cockroachdb/errors"
	cfg "github.com/neatplex/nightell-core/internal/config"
	"github.com/neatplex/nightell-core/internal/logger"
	"io"
)

type LoggerProxy struct {
	e *logger.Logger
}

func (l *LoggerProxy) Logf(level logging.Classification, format string, a ...interface{}) {
	if level == logging.Debug {
		l.e.Debug(fmt.Sprintf(format, a...))
	} else {
		l.e.Warn(fmt.Sprintf(format, a...))
	}
}

type S3 struct {
	config *cfg.Config
	l      *logger.Logger
	client *s3.Client
}

func (s *S3) Init() error {
	cs := []func(options *config.LoadOptions) error{
		config.WithRegion(s.config.S3.Region),
		config.WithLogger(&LoggerProxy{e: s.l}),
		config.WithClientLogMode(aws.LogRetries | aws.LogRequest | aws.LogResponse),
	}

	if !s.config.S3.RoleUsed {
		credentialsCache := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			s.config.S3.AccessKey,
			s.config.S3.SecretKey,
			"",
		))
		cs = append(cs, config.WithCredentialsProvider(credentialsCache), config.WithRegion(s.config.S3.Region))
	}

	c, err := config.LoadDefaultConfig(context.TODO(), cs...)
	if err != nil {
		return errors.Wrap(err, "s3: cannot load s3 config")
	} else {
		s.l.Debug("s3: connection established")
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
		return nil, errors.Wrap(err, "s3: cannot download file")
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	body, err := io.ReadAll(r.Body)
	return body, errors.Wrap(err, "s3: cannot read downloaded file")
}

func (s *S3) Put(path string, body io.Reader) error {
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &s.config.S3.Bucket,
		Key:    &path,
		Body:   body,
	})
	return errors.Wrap(err, "s3: cannot upload")
}

func (s *S3) Delete(path string) error {
	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &s.config.S3.Bucket,
		Key:    &path,
	})
	return errors.Wrap(err, "s3: cannot delete")
}

func New(c *cfg.Config, l *logger.Logger) *S3 {
	return &S3{config: c, l: l}
}
