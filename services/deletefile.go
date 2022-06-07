package services

import (
	"context"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

// grabbed this from: https://aws.github.io/aws-sdk-go-v2/docs/code-examples/s3/deleteobject/

type S3DeleteObjectAPI interface {
	DeleteObject(ctx context.Context,
		params *s3.DeleteObjectInput,
		optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

func DeleteItem(c context.Context, api S3DeleteObjectAPI, input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return api.DeleteObject(c, input)
}

func DeleteS3Object(fileObject, awsRegion, s3Bucket string) (err error) {

	acfg, err := awscfg.LoadDefaultConfig(context.TODO(),
		awscfg.WithRegion(awsRegion),
	)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(acfg)

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(fileObject),
	}

	_, err = DeleteItem(context.TODO(), client, input)
	if err != nil {
		return err
	}

	return nil
}
