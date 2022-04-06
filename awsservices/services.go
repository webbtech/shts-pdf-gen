package awsservices

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// PutFile function
func PutFile(buf *bytes.Buffer, fileObject, awsRegion, s3Bucket string) (location string, err error) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		return "", err
	}

	uploader := s3manager.NewUploader(sess)
	res, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:             aws.String(s3Bucket),
		Key:                aws.String(fileObject),
		Body:               buf,
		ContentType:        aws.String("application/pdf"),
		ContentDisposition: aws.String("attachment"),
	})
	if err != nil {
		return "", err
	}

	return string(res.Location), nil
}
