package aws

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	AwsS3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//S3Factory - S3 helper type
type S3Factory struct {
	Bucket string
	Region string
}

//NewS3 - S3 bucket initializer
func NewS3() *S3Factory {
	return &S3Factory{
		Bucket: "sharrans-bucket",
		Region: "us-east-1",
	}
}

//UploadFile - upload files to s3 at the given path. Bucket and region are pre-configured
func (s *S3Factory) UploadFile(formDataStr, filename string) (string, error) {
	reader := ioutil.NopCloser(bytes.NewBufferString(formDataStr))
	defer reader.Close()

	sess, err := NewSession(s.Region)
	uploader := s3manager.NewUploader(sess)
	res, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(filename),
		Body:   reader,
	})
	if err != nil {
		fmt.Println("s3Factory :", err.Error())
		return "", err
	}

	fmt.Printf("Successfully uploaded %q to %q\n", filename, s.Bucket)
	return res.Location, nil
}

//GetPreSignedURL - get signed url for the provided file path
func (s *S3Factory) GetPreSignedURL(filename string) (string, error) {
	sess, err := NewSession(s.Region)
	if err != nil {
		return "", err
	}

	req, _ := AwsS3.New(sess).GetObjectRequest(&AwsS3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(filename),
	})

	urlStr, err := req.Presign(15 * time.Minute)

	if err != nil {
		return "", err
	}

	return urlStr, nil
}
