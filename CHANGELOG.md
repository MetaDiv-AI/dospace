# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.5] - 2025-02-16

### Fixed

- `GetPublicURL` virtual-hosted-style: bucket in subdomain only (`https://{bucket}.{region}.digitaloceanspaces.com/{key}`)

## [1.0.4] - 2025-02-16

### Fixed

- `GetPublicURL` now includes bucket in the path for correct DigitalOcean Spaces URLs (`https://{bucket}.{region}.digitaloceanspaces.com/{bucket}/{key}`)

## [1.0.3] - 2025-02-16

### Fixed

- `GetPublicURL` virtual-hosted-style URLs: bucket in subdomain only, not in path (`https://{bucket}.{region}.digitaloceanspaces.com/{key}`)

## [1.0.2] - 2025-02-16

### Fixed

- `GetPublicURL` now includes bucket in the path for correct virtual-hosted-style URLs

## [1.0.1] - 2025-02-15

### Changed

- `GetPublicURL` now uses virtual-hosted-style URLs (`https://{bucket}.{region}.digitaloceanspaces.com/{key}`) instead of path-style, so the bucket is in the subdomain rather than the path

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
