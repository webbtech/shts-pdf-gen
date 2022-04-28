package services

import (
	"bytes"
	"context"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

// S3PutObjectAPI defines the interface for the PutObject function.
// We use this interface to test the function using a mocked service.
type S3PutObjectAPI interface {
	PutObject(ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

// PutFile uploads a file to an Amazon Simple Storage Service (Amazon S3) bucket
// Inputs:
//     c is the context of the method call, which includes the AWS Region
//     api is the interface that defines the method call
//     input defines the input arguments to the service call.
// Output:
//     If success, a PutObjectOutput object containing the result of the service call and nil
//     Otherwise, nil and an error from the call to PutObject
func PutFile(c context.Context, api S3PutObjectAPI, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return api.PutObject(c, input)
}

func UploadS3Object(buf *bytes.Buffer, fileObject, awsRegion, s3Bucket string) (err error) {

	acfg, err := awscfg.LoadDefaultConfig(context.TODO(),
		awscfg.WithRegion(awsRegion),
	)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(acfg)

	input := &s3.PutObjectInput{
		Bucket:             aws.String(s3Bucket),
		Key:                aws.String(fileObject),
		Body:               buf,
		ContentType:        aws.String("application/pdf"),
		ContentDisposition: aws.String("attachment"),
	}

	_, err = PutFile(context.TODO(), client, input)
	if err != nil {
		return err
	}

	return nil
}
