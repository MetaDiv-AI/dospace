// Package do_spaces provides a client for interacting with DigitalOcean Spaces (S3-compatible storage).
package do_spaces

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// NewClient creates a new Client for DigitalOcean Spaces.
// endpoint should be like "https://sgp1.digitaloceanspaces.com"
// region should be like "sgp1"
// accessKey and secretKey are the DO Spaces credentials
func NewClient(endpoint string, region string, bucket string, accessKey string, secretKey string) (*Client, error) {
	if endpoint == "" || bucket == "" || accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("endpoint, bucket, accessKey, and secretKey are required")
	}
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, reg string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           endpoint,
			SigningRegion: region,
		}, nil
	})

	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		awsconfig.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // DO Spaces requires path-style addressing
	})

	return &Client{
		S3Client: s3Client,
		Bucket:   bucket,
		Endpoint: endpoint,
	}, nil
}

// Client represents a DigitalOcean Spaces client with a specific bucket.
type Client struct {
	S3Client *s3.Client
	Bucket   string
	Endpoint string // e.g. "https://sgp1.digitaloceanspaces.com"
}

// Upload uploads the given content as an object with the specified name to the bucket.
// contentType is optional - if provided, it will set Content-Type and Content-Disposition headers.
// It returns an error if the upload fails.
func (c *Client) Upload(ctx context.Context, objectName string, content []byte, contentType ...string) error {
	return c.UploadFromReader(ctx, objectName, bytes.NewReader(content), contentType...)
}

// UploadFromReader uploads content from an io.Reader to the bucket.
// Use this for large files to avoid loading the entire content into memory.
func (c *Client) UploadFromReader(ctx context.Context, objectName string, r io.Reader, contentType ...string) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(objectName),
		Body:   r,
		ACL:    types.ObjectCannedACLPublicRead,
	}

	if len(contentType) > 0 && contentType[0] != "" {
		input.ContentType = aws.String(contentType[0])
		if shouldBeInline(contentType[0]) {
			input.ContentDisposition = aws.String("inline")
		}
	}

	_, err := c.S3Client.PutObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to upload object: %w", err)
	}
	return nil
}

// shouldBeInline determines if a content type should be displayed inline in browsers
func shouldBeInline(contentType string) bool {
	contentType = strings.ToLower(contentType)
	return strings.HasPrefix(contentType, "image/") ||
		strings.HasPrefix(contentType, "video/") ||
		strings.HasPrefix(contentType, "audio/") ||
		contentType == "application/pdf"
}

// Download retrieves the object with the specified name from the bucket and returns its content as a byte slice.
func (c *Client) Download(ctx context.Context, objectName string) ([]byte, error) {
	result, err := c.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(objectName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

// Delete removes the object with the specified name from the bucket.
func (c *Client) Delete(ctx context.Context, objectName string) error {
	_, err := c.S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(objectName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

// Exists checks if an object exists in the bucket.
// Returns true if the object exists, false if it does not (404), or an error for other failures.
func (c *Client) Exists(ctx context.Context, objectName string) (bool, error) {
	_, err := c.S3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(objectName),
	})
	if err != nil {
		var notFound *types.NotFound
		if errors.As(err, &notFound) {
			return false, nil
		}
		// Fallback for S3-compatible services that may not use types.NotFound
		var apiErr interface{ ErrorCode() string }
		if errors.As(err, &apiErr) {
			code := apiErr.ErrorCode()
			if code == "NotFound" || code == "NoSuchKey" {
				return false, nil
			}
		}
		return false, fmt.Errorf("failed to head object: %w", err)
	}
	return true, nil
}

// GetPublicURL returns the public URL for an object in DigitalOcean Spaces.
// Uses virtual-hosted-style: https://{bucket}.{region}.digitaloceanspaces.com/{objectKey}
func (c *Client) GetPublicURL(objectName string) string {
	baseURL := strings.TrimPrefix(strings.TrimPrefix(c.Endpoint, "https://"), "http://")
	return fmt.Sprintf("https://%s/%s/%s", baseURL, c.Bucket, objectName)
}
