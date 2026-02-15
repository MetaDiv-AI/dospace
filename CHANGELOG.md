# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-02-15

### Added

- Initial release
- `NewClient` - Create a client for DigitalOcean Spaces with endpoint, region, bucket, and credentials
- `Upload` - Upload bytes to the bucket with optional Content-Type
- `UploadFromReader` - Stream upload from `io.Reader` for large files
- `Download` - Download object as `[]byte`
- `List` - List object keys with optional prefix filter (handles pagination)
- `Exists` - Check if an object exists (with S3-compatible error fallback)
- `GetPublicURL` - Get public URL for an object
- `Delete` - Delete an object
- Content-Disposition: inline for previewable types (images, video, audio, PDF)
- Input validation for `NewClient` (endpoint, bucket, accessKey, secretKey required)
- Context support for cancellation and timeouts on all operations
