package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3api struct {
	client       *s3.S3
	nativeBucket string
}

func newS3API(bucket string) (*s3api, error) {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("Failed to create session, ", err)
		return nil, err
	}

	return &s3api{
		s3.New(sess),
		bucket,
	}, nil
}

func (s3api *s3api) Ids(collection string) ([]string, error) {
	ids := make([]string, 0)
	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(s3api.nativeBucket),
		// Prefix: aws.String(collection),
	}

	err := s3api.client.ListObjectsV2Pages(params,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			for _, p := range page.Contents {
				if *p.Size > int64(0) {
					ids = append(ids, *p.Key)
				}
			}
			return lastPage
		})
	if err != nil {
		return nil, err
	}
	return ids, nil
}
