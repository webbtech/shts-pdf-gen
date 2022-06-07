package services

import (
	"bytes"
	"context"
	"os"
	"testing"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/webbtech/shts-pdf-gen/config"
)

type S3PutObjectAPIImp struct{}

func (dt S3PutObjectAPIImp) PutObject(ctx context.Context,
	params *s3.PutObjectInput,
	optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {

	output := &s3.PutObjectOutput{}

	return output, nil
}

// These tests aren't exactly a unit test as it depends on our Config object,
// but it seems like a minor rule to break?...

func TestPutObject(t *testing.T) {
	api := &S3PutObjectAPIImp{}

	getConfig(t)

	keyObj := "estimate/est-99.pdf"

	filename := "test.txt"
	file, _ := os.Open(filename)

	defer file.Close()

	input := &s3.PutObjectInput{
		Bucket: aws.String(cfg.S3Bucket),
		Key:    aws.String(keyObj),
		Body:   file,
	}

	resp, err := PutFile(context.TODO(), *api, input)
	if err != nil {
		t.Fatalf("Expected null error received: %s", err)
	}

	t.Logf("resp: %+v\n", resp)
}

func TestIntegPutObject(t *testing.T) {

	getConfig(t)

	keyObj := "estimate/est-99.pdf"
	buf := bytes.NewBuffer(nil)

	err := UploadS3Object(buf, keyObj, cfg.AwsRegion, cfg.S3Bucket)
	if err != nil {
		t.Fatalf("Expected null error received: %s", err)
	}

	// teardown
	deleteObject(t, keyObj)
}

// ================================== Helpers ==========================================

var cfg *config.Config

func getConfig(t *testing.T) {
	t.Helper()

	cfg = &config.Config{}
	cfg.Init()
}

func deleteObject(t *testing.T, objectName string) {

	t.Helper()
	getConfig(t)

	acfg, err := awscfg.LoadDefaultConfig(context.TODO(),
		awscfg.WithRegion(cfg.AwsRegion),
	)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	client := s3.NewFromConfig(acfg)

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(cfg.S3Bucket),
		Key:    aws.String(objectName),
	}

	_, err = client.DeleteObject(context.TODO(), input)
	if err != nil {
		t.Fatalf("Got an error deleting item: %s", err)
	}
}
