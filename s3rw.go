package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3api struct {
	client       *s3.S3
	nativeBucket string
}

type s3resource struct {
	uuid        string
	content     []byte
	contentType string
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

func (s3a *s3api) Ids(collection string) ([]string, error) {
	var ids []string
	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(s3a.nativeBucket),
		Prefix: aws.String(collection),
	}

	err := s3a.client.ListObjectsV2Pages(params,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			for _, p := range page.Contents {
				if *p.Size > int64(0) {
					ids = append(ids, strings.Join(strings.Split(*p.Key, "/")[1:], "-"))
				}
			}
			return lastPage
		})
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (s3a *s3api) Read(collection, uuid string) (s3resource, error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(s3a.nativeBucket),
		Key:    aws.String(uuidToS3Key(collection, uuid)),
	}

	resp, err := s3a.client.GetObject(params)
	content, _ := ioutil.ReadAll(resp.Body)

	return s3resource{
		uuid:        uuid,
		content:     content,
		contentType: *resp.ContentType,
	}, err
}

func (s3a *s3api) Write(collection string, resource s3resource) error {
	body := bytes.NewReader(resource.content)
	s3uuid := uuidToS3Key(collection, resource.uuid)
	params := &s3.PutObjectInput{
		Bucket:      aws.String(s3a.nativeBucket),
		Key:         &s3uuid,
		Body:        body,
		ContentType: &resource.contentType,
	}
	_, err := s3a.client.PutObject(params)

	return err
}

func (s3a *s3api) Delete(collection, uuid string) error {
	s3uuid := uuidToS3Key(collection, uuid)
	params := &s3.DeleteObjectInput{
		Bucket: aws.String(s3a.nativeBucket),
		Key:    &s3uuid,
	}

	_, err := s3a.client.DeleteObject(params)
	return err
}

func uuidToS3Key(collection, uuid string) string {
	return collection + "/" + strings.Replace(uuid, "-", "/", 4)
}
