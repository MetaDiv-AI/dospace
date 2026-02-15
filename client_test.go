package do_spaces

import (
	"testing"
)

func TestShouldBeInline(t *testing.T) {
	tests := []struct {
		contentType string
		want       bool
	}{
		{"image/png", true},
		{"image/jpeg", true},
		{"IMAGE/PNG", true},
		{"video/mp4", true},
		{"audio/mpeg", true},
		{"application/pdf", true},
		{"text/plain", false},
		{"application/octet-stream", false},
		{"application/json", false},
		{"", false},
	}
	for _, tt := range tests {
		got := shouldBeInline(tt.contentType)
		if got != tt.want {
			t.Errorf("shouldBeInline(%q) = %v, want %v", tt.contentType, got, tt.want)
		}
	}
}

func TestGetPublicURL(t *testing.T) {
	tests := []struct {
		endpoint   string
		bucket     string
		objectName string
		want       string
	}{
		{
			endpoint:   "https://sgp1.digitaloceanspaces.com",
			bucket:     "my-bucket",
			objectName: "path/to/file.txt",
			want:       "https://sgp1.digitaloceanspaces.com/my-bucket/path/to/file.txt",
		},
		{
			endpoint:   "https://nyc3.digitaloceanspaces.com",
			bucket:     "assets",
			objectName: "image.png",
			want:       "https://nyc3.digitaloceanspaces.com/assets/image.png",
		},
		{
			endpoint:   "http://sgp1.digitaloceanspaces.com",
			bucket:     "bucket",
			objectName: "file",
			want:       "https://sgp1.digitaloceanspaces.com/bucket/file",
		},
	}
	for _, tt := range tests {
		client := &Client{Bucket: tt.bucket, Endpoint: tt.endpoint}
		got := client.GetPublicURL(tt.objectName)
		if got != tt.want {
			t.Errorf("GetPublicURL(%q) = %q, want %q", tt.objectName, got, tt.want)
		}
	}
}

func TestNewClient_Validation(t *testing.T) {
	tests := []struct {
		endpoint, region, bucket, accessKey, secretKey string
		wantErr                                      bool
	}{
		{"", "sgp1", "bucket", "key", "secret", true},
		{"https://sgp1.digitaloceanspaces.com", "sgp1", "", "key", "secret", true},
		{"https://sgp1.digitaloceanspaces.com", "sgp1", "bucket", "", "secret", true},
		{"https://sgp1.digitaloceanspaces.com", "sgp1", "bucket", "key", "", true},
		{"https://sgp1.digitaloceanspaces.com", "sgp1", "bucket", "key", "secret", false},
	}
	for _, tt := range tests {
		_, err := NewClient(tt.endpoint, tt.region, tt.bucket, tt.accessKey, tt.secretKey)
		if (err != nil) != tt.wantErr {
			t.Errorf("NewClient(...) err = %v, wantErr %v", err, tt.wantErr)
		}
	}
}

func TestNewClient_ReturnsClient(t *testing.T) {
	client, err := NewClient(
		"https://sgp1.digitaloceanspaces.com",
		"sgp1",
		"test-bucket",
		"test-key",
		"test-secret",
	)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	if client == nil {
		t.Fatal("NewClient returned nil client")
	}
	if client.Bucket != "test-bucket" {
		t.Errorf("Bucket = %q, want test-bucket", client.Bucket)
	}
	if client.Endpoint != "https://sgp1.digitaloceanspaces.com" {
		t.Errorf("Endpoint = %q, want https://sgp1.digitaloceanspaces.com", client.Endpoint)
	}
}

func TestUploadDownloadRoundTrip(t *testing.T) {
	// Skip if no credentials - this would be an integration test
	// that requires DO_SPACES_* env vars or similar
	t.Skip("integration test - requires DO Spaces credentials")
}
