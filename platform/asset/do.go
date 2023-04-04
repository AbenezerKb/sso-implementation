package asset

import (
	"context"
	"fmt"
	"mime"
	"mime/multipart"
	"path"

	"sso/platform"
	"sso/platform/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.uber.org/zap"
)

type digitalOceanAsset struct {
	log                      logger.Logger
	key, secret, url, bucket string
	s3Client                 *s3.S3
}

func InitDigitalOceanAsset(log logger.Logger, key, secret, url, bucket string) platform.Asset {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:         aws.String(url),
		S3ForcePathStyle: aws.Bool(false),
		Region:           aws.String("us-east-1"),
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		log.Fatal(context.Background(), "could not create a session for digital ocean s3 client", zap.Error(err))
	}

	do := &digitalOceanAsset{
		log:      log,
		key:      key,
		secret:   secret,
		url:      url,
		bucket:   bucket,
		s3Client: s3.New(newSession),
	}

	spaces, err := do.s3Client.ListBuckets(nil)
	if err != nil {
		log.Fatal(context.Background(), "could not list buckets for digital ocean s3 client", zap.Error(err))
	}

	var found bool

	for _, b := range spaces.Buckets {
		fmt.Println(aws.StringValue(b.Name))
		if aws.StringValue(b.Name) == bucket {
			found = true

			break
		}
	}

	if !found {
		if err := do.createBucket(); err != nil {
			log.Fatal(context.Background(), "could not create bucket for digital ocean s3 client", zap.Error(err))
		}
	}

	return do
}

func (f *digitalOceanAsset) createBucket() error {
	_, err := f.s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(f.bucket),
	})
	if err != nil {
		return err
	}

	return nil
}

func (f *digitalOceanAsset) SaveAsset(ctx context.Context, asset multipart.File, dst string) error {
	object := s3.PutObjectInput{
		Bucket:      aws.String(f.bucket),
		Key:         aws.String(dst),
		Body:        asset,
		ACL:         aws.String("private"),
		ContentType: aws.String(mime.TypeByExtension(path.Ext(dst))),
	}

	_, err := f.s3Client.PutObjectWithContext(ctx, &object)
	if err != nil {
		return err
	}

	return nil
}
