package cmhp_s3

import (
	"bytes"
	"io"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/maldan/go-cmhp/cmhp_file"
)

var s3Client *s3.S3
var config struct {
	SPACES_KEY      string `json:"SPACES_KEY"`
	SPACES_SECRET   string `json:"SPACES_SECRET"`
	SPACES_ENDPOINT string `json:"SPACES_ENDPOINT"`
	SPACES_BUCKET   string `json:"SPACES_BUCKET"`
}

type WriteArgs struct {
	Path        string
	InputData   []byte
	Visibility  string
	ContentType string
	MetaData    map[string]string
}

type S3File struct {
	Path         string
	Size         int64
	LastModified time.Time
}

func Start(path string) {
	cmhp_file.ReadJSON(path, &config)

	newSession, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(config.SPACES_KEY, config.SPACES_SECRET, ""),
		Endpoint:    aws.String(config.SPACES_ENDPOINT),
		Region:      aws.String("us-east-1"),
	})
	if err != nil {
		log.Fatal(err)
	}
	s3Client = s3.New(newSession)
	log.Println("S3 is ready")
}

func List(path string) []S3File {
	list, _ := s3Client.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(config.SPACES_BUCKET),
		Prefix: aws.String(path),
	})
	out := make([]S3File, 0)
	for _, o := range list.Contents {

		out = append(out, S3File{
			Path:         *o.Key,
			Size:         *o.Size,
			LastModified: *o.LastModified,
		})
	}
	return out
}

func Write(args WriteArgs) (string, error) {
	// Url of endpoint
	url := strings.Replace(
		config.SPACES_ENDPOINT,
		"https://", "https://"+config.SPACES_BUCKET+".",
		1,
	)

	// Remove https://...
	path := pureS3Path(args.Path)

	// Write object
	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(config.SPACES_BUCKET),
		Key:         aws.String(path),
		Body:        bytes.NewReader(args.InputData),
		ACL:         aws.String(args.Visibility),
		ContentType: aws.String(args.ContentType),
		Metadata:    aws.StringMap(args.MetaData),
	})
	if err != nil {
		return "", err
	}

	return url + "/" + path, nil
}

func Read(path string) ([]byte, error) {
	// Remove https://...
	path = pureS3Path(path)

	// Download file
	result, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(config.SPACES_BUCKET),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}

	// Read and decompress
	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func pureS3Path(path string) string {
	// Url of endpoint
	url := strings.Replace(
		config.SPACES_ENDPOINT,
		"https://", "https://"+config.SPACES_BUCKET+".",
		1,
	)
	path = strings.Replace(path, url+"/", "", 1)
	return path
}
