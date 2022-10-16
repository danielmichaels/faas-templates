package handlers

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type DataStore interface {
	// Download a file from s3 passing the filepath in s3, and the filepath to
	// download that file to locally
	Download(from string, to string) error
	// List lists all the files within the bucket
	List() ([]string, error)
}
type Store struct {
	Data DataStore
}
type B2Client struct {
	bucketName string
	s3Client   *s3.S3
}

// getSecret retrieves the secret from openfaas and makes it available for use.
func getSecret(secretName string) ([]byte, error) {
	secret, err := ioutil.ReadFile(fmt.Sprintf("/var/openfaas/secrets/%s", secretName))
	if err != nil {
		return nil, err
	}

	s := strings.TrimSpace(string(secret))
	return []byte(s), nil
}

func NewDataStore() (*Store, error) {
	b2AppKey, err := getSecret("b2AppKey")
	if err != nil {
		log.Fatalln("no b2AppKey found")
	}
	b2KeyId, err := getSecret("b2KeyID")
	if err != nil {
		log.Fatalln("no b2KeyID found")
	}
	endpoint, err := getSecret("b2Server")
	if err != nil {
		log.Fatalln("no b2Server found")
	}
	b2Bucket, err := getSecret("b2Bucket")
	if err != nil {
		log.Fatalln("no b2Bucket found")
	}

	cfg := &aws.Config{
		Endpoint:    aws.String(string(endpoint)),
		Region:      aws.String("us-east-002"),
		Credentials: credentials.NewStaticCredentials(string(b2KeyId), string(b2AppKey), ""),
	}

	s, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}

	cl := s3.New(s)

	return &Store{
		Data: B2Client{
			bucketName: string(b2Bucket),
			s3Client:   cl,
		},
	}, nil
}

// Download downloads a file from s3 by giving the filepath in s3 and filepath to download
// the file to.
func (b B2Client) Download(filename string, to string) error {
	file, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("b2.download error: failed creating destination file: '%s'", err)
	}
	defer file.Close()

	dl := s3manager.NewDownloaderWithClient(b.s3Client)
	_, err = dl.Download(file, &s3.GetObjectInput{
		Bucket: &b.bucketName,
		Key:    &filename,
	})
	if err != nil {
		_ = os.Remove(to)
		return fmt.Errorf("b2.download error: failed to download file '%s'", err)
	}
	return nil
}

func (b B2Client) List() ([]string, error) {
	objects, err := b.s3Client.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: &b.bucketName})
	if err != nil {
		return nil, fmt.Errorf("b2.list error: failed to list objects '%s'", err)
	}

	result := make([]string, 0, len(objects.Contents))
	for _, k := range objects.Contents {
		result = append(result, *k.Key)
	}
	return result, nil
}
