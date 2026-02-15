# dospace

A Go client for [DigitalOcean Spaces](https://docs.digitalocean.com/products/spaces/) (S3-compatible object storage).

## Installation

```bash
go get github.com/MetaDiv-AI/dospace
```

## Usage

```go
package main

import (
    "context"
    "log"

    "github.com/MetaDiv-AI/dospace"
)

func main() {
    client, err := do_spaces.NewClient(
        "https://sgp1.digitaloceanspaces.com", // endpoint
        "sgp1",                                 // region
        "my-bucket",                            // bucket name
        "your-access-key",
        "your-secret-key",
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Upload
    if err := client.Upload(ctx, "path/to/file.txt", []byte("hello"), "text/plain"); err != nil {
        log.Fatal(err)
    }

    // Upload from reader (for large files)
    // client.UploadFromReader(ctx, "large.bin", file, "application/octet-stream")

    // Download
    data, err := client.Download(ctx, "path/to/file.txt")
    if err != nil {
        log.Fatal(err)
    }

    // List objects (optional prefix)
    keys, err := client.List(ctx, "path/")
    if err != nil {
        log.Fatal(err)
    }

    // Check if object exists
    exists, err := client.Exists(ctx, "path/to/file.txt")
    if err != nil {
        log.Fatal(err)
    }

    // Get public URL
    url := client.GetPublicURL("path/to/file.txt")

    // Delete
    if err := client.Delete(ctx, "path/to/file.txt"); err != nil {
        log.Fatal(err)
    }
}
```

## API

| Method | Description |
|--------|-------------|
| `NewClient(endpoint, region, bucket, accessKey, secretKey)` | Create a new client |
| `Upload(ctx, objectName, content, contentType...)` | Upload bytes to the bucket |
| `UploadFromReader(ctx, objectName, r, contentType...)` | Stream upload from `io.Reader` |
| `Download(ctx, objectName)` | Download object as `[]byte` |
| `List(ctx, prefix)` | List object keys, optionally filtered by prefix |
| `Exists(ctx, objectName)` | Check if object exists |
| `GetPublicURL(objectName)` | Get public URL for an object |
| `Delete(ctx, objectName)` | Delete an object |

## Credentials

Create access keys from the [DigitalOcean Spaces Access Keys](https://cloud.digitalocean.com/spaces/access_keys) page. The endpoint format is `https://<region>.digitaloceanspaces.com` (e.g. `sgp1`, `nyc3`).

## License

See repository.
