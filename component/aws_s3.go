package component

import (
	"bytes"
	"context"
	"flag"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"net/http"

	sctx "github.com/viettranx/service-context"
)

type UploadProvider interface {
	SaveFileUploaded(ctx context.Context, data []byte, dst string) error
}

type s3Provider struct {
	id         string
	bucketName string
	region     string
	domain     string
	client    *s3.Client
}

func (p *s3Provider) ID() string { return p.id }

func (p *s3Provider) InitFlags() {
	// S3 env variable are not consider as global variable in AWS so we can not get that from the config package
	// We will use flagenv package to parse flag "aws-s3-bucket" to custom env variable and assign it to p.bucketNam instead 

	flag.StringVar(
		&p.bucketName,
		"aws-s3-bucket",
		"",
		"ASW S3 bucket name",
	)

	flag.StringVar(
		&p.region,
		"aws-region",
		"ap-southeast-1",
		"ASW region name",
	)

	flag.StringVar(
		&p.domain,
		"cdn-domain",
		"",
		"ASW S3 cdn domain name",
	)
}

func NewAWSS3Provider(id string) *s3Provider {
	return &s3Provider{id: id}
}

func (p *s3Provider) Activate(_ sctx.ServiceContext) error {
	// Create new AWS client
	/*
		// Use aws.Config to custom modify config variable
		cfg := aws.Config{
			Region: p.region,
			....
		}
	*/
	// Use LoadDefaultConfig to load default aws config (access, secret key, region (not include S3))
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	p.client = s3.NewFromConfig(cfg)
	return nil
}

func (p *s3Provider) Stop() error {
	return nil
}

func (p *s3Provider) SaveFileUploaded(ctx context.Context, data []byte, dst string) error {
	fileBytes := bytes.NewReader(data)
	fileType := http.DetectContentType(data)

	_, err := p.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(p.bucketName),
		Key:         aws.String(dst),
		ACL:         types.ObjectCannedACLPrivate,
		ContentType: aws.String(fileType),
		Body:        fileBytes,
	})

	if err != nil {
		return err
	}

	return nil
}

// func (p *s3Provider) GetUploadPresignedURL(ctx context.Context) string {
// 	req, _ := s3.New(p.session).PutObjectRequest(&s3.PutObjectInput{
// 		Bucket: aws.String(p.bucketName),
// 		Key:    aws.String(fmt.Sprintf("img/%d", time.Now().UnixNano())),
// 		ACL:    aws.String("private"),
// 	})
// 	//
// 	url, _ := req.Presign(time.Second * 60)

// 	return url
// }

func (p *s3Provider) GetDomain() string { return p.domain }
func (*s3Provider) GetName() string     { return "aws_s3" }