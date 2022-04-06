package awsservices

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/webbtech/shts-pdf-gen/config"
)

// PutFile function
func PutFile(fn string, buf *bytes.Buffer, cfg *config.Config) (location string, err error) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.AwsRegion),
	})
	if err != nil {
		return "", err
	}

	uploader := s3manager.NewUploader(sess)
	res, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:             aws.String(cfg.S3Bucket),
		Key:                aws.String(fn),
		Body:               buf,
		ContentType:        aws.String("application/pdf"),
		ContentDisposition: aws.String("attachment"),
	})
	if err != nil {
		return "", err
	}

	return string(res.Location), nil
}
