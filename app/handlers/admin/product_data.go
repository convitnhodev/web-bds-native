package web

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"time"
)

type Furniture struct {
	ID    int
	Value string
}

type BrowserS3Policy struct {
	Expiration string        `json:"expiration"`
	Conditions []interface{} `json:"conditions"`
}

var furniture = []Furniture{
	{
		ID:    1,
		Value: "Nội thất đầy đủ",
	},
	{
		ID:    2,
		Value: "Nội thất cơ bản",
	},
	{
		ID:    3,
		Value: "Không nội thất",
	},
	{
		ID:    4,
		Value: "Nhà thô",
	},
	{
		ID:    5,
		Value: "Khác",
	},
}

var s3Bucket = "d1e2e3i4n5"
var s3accessKeyID = "002026794d2122d000000000c"
var s3Secret = "K0022S/oKffHhTFf12PpuVQjim2O0Xw"

// var s3Endpoint = "//d1e2e3i4n5.s3.us-west-002.backblazeb2.com"
// var s3Endpoint = "https://s3.us-west-002.backblazeb2.com/"
var s3Endpoint = "https://s3.us-west-002.backblazeb2.com/d1e2e3i4n5"

func GetS3Policy() string {
	conditions := make([]interface{}, 0)
	conditions = append(conditions, map[string]string{"bucket": s3Bucket})
	conditions = append(conditions, map[string]string{"acl": "public-read"})
	conditions = append(conditions, []string{"starts-with", "$key", ""})
	conditions = append(conditions, []string{"starts-with", "$Content-Type", "image/"})
	conditions = append(conditions, []string{"starts-with", "$x-amz-meta-tag", ""})

	conditions = append(conditions, map[string]string{"x-amz-meta-uuid": "14365123651274"})
	conditions = append(conditions, map[string]string{"x-amz-server-side-encryption": "AES256"})

	conditions = append(conditions, map[string]string{"x-amz-credential": "002026794d2122d000000000c/20221023/us-west-002/s3/aws4_request"})
	conditions = append(conditions, map[string]string{"x-amz-algorithm": "AWS4-HMAC-SHA256"})
	conditions = append(conditions, map[string]string{"x-amz-date": "20221023T000000Z"})

	conditions = append(conditions, []string{"starts-with", "$name", ""})
	conditions = append(conditions, []string{"starts-with", "$Filename", ""})

	policy := BrowserS3Policy{
		Expiration: time.Now().Add(15 * time.Minute).Format("2006-01-02T15:04:05.999Z"),
		Conditions: conditions,
	}
	jp, _ := json.Marshal(&policy)
	return string(jp)
}

func GetSignature() string {
	p := GetS3Policy()
	return hashSignature(p, s3Secret)
}

func hashSignature(input string, key string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(input))
	s := hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func GetS3BrowserSecret() (accessKeyID string, policy string, signature string) {
	return s3accessKeyID, GetS3Policy(), GetSignature()
}

func GetAWSS3Policy() string {
	cfg := aws.NewConfig()
	cfg.Region = "us-west-002"
	cfg.Credentials = aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     "002026794d2122d000000000c",
			SecretAccessKey: "K0022S/oKffHhTFf12PpuVQjim2O0Xw",
		}, nil
	})
	cfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           "https://s3.us-west-002.backblazeb2.com",
			Source:        aws.EndpointSourceCustom,
			SigningName:   "backblazeb2",
			SigningRegion: region,
		}, nil
	})
	//return ""

	client := s3.NewFromConfig(*cfg)
	presignClient := s3.NewPresignClient(client)
	presignDuration := func(po *s3.PresignOptions) {
		po.Expires = 5 * time.Minute
	}
	//id := xid.New().String()
	//
	presignParams := &s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String("test/ahihi.jpg"),
	}

	presignResult, err := presignClient.PresignPutObject(context.Background(), presignParams, presignDuration)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("Presigned URL For object: %s\n", s.URL)
	fmt.Printf("Signed: %+v\n", presignResult)
	return ""
}
