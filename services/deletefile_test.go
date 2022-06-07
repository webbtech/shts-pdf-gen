package services

import (
	"bytes"
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

type S3DeleteObjectAPIImp struct{}

var keyObj = "estimate/est-99.pdf"

func (dt S3DeleteObjectAPIImp) DeleteObject(ctx context.Context,
	params *s3.DeleteObjectInput,
	optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {

	output := &s3.DeleteObjectOutput{}

	return output, nil
}

// These tests aren't exactly a unit test as it depends on our Config object,
// but it seems like a minor rule to break?...

func TestDeleteObject(t *testing.T) {

	api := &S3DeleteObjectAPIImp{}

	getConfig(t)

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(cfg.S3Bucket),
		Key:    aws.String(keyObj),
	}

	resp, err := DeleteItem(context.TODO(), *api, input)
	if err != nil {
		t.Fatalf("Expected null error received: %s", err)
	}

	t.Logf("resp: %+v\n", resp)
}

func TestIntegDeleteObject(t *testing.T) {

	getConfig(t)

	// setup
	putObject(t)

	if err := DeleteS3Object(keyObj, cfg.AwsRegion, cfg.S3Bucket); err != nil {
		t.Fatalf("Expected null error received: %s", err)
	}
}

// ================================== Helpers ==========================================
func putObject(t *testing.T) {

	getConfig(t)

	buf := bytes.NewBuffer(nil)

	err := UploadS3Object(buf, keyObj, cfg.AwsRegion, cfg.S3Bucket)
	if err != nil {
		t.Fatalf("Expected null error received: %s", err)
	}
}
