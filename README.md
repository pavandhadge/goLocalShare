# Secure File Server

A secure HTTP server for sharing files and directories with authentication and rate limiting.

## Features

- **Secure file sharing**: Share single files or entire directories
- **Authentication**: Time-limited access tokens (valid for 1 hour)
- **Rate limiting**: Protection against brute force attacks
- **Security headers**: CSP, XSS protection, no-sniff, etc.
- **Path sanitization**: Protection against directory traversal attacks
- **Two modes**:
  - Single file mode (direct download link)
  - Directory mode (browsable interface)

## Installation

1. Ensure you have Go installed (version 1.16+ recommended)
2. Clone or download the repository
3. Build the binary:
   ```sh
   go build -o server
   ```

## Usage

### For a single file:
```sh
./server /path/to/your/file.ext
```

### For a directory:
```sh
./server --dir /path/to/your/directory
```

The server will start on `http://localhost:8080`

## Access Control

- Upon visiting the homepage (`/`), you'll receive an access token
- This token is valid for 1 hour
- All subsequent requests require this token either:
  - As a URL parameter (`?token=...`)
  - Or in the `X-Auth-Token` header

## Security Features

- **Rate limiting**: 100 requests per minute per IP
- **Brute force protection**: 2-second delay after failed auth attempts
- **Secure headers**: CSP, XSS protection, no-sniff, etc.
- **Path sanitization**: Prevents directory traversal attacks
- **Symlink protection**: Blocks symlinks that point outside the base directory

## API Endpoints

| Endpoint       | Description                                                                 |
|----------------|-----------------------------------------------------------------------------|
| `/`            | Homepage - displays access token and appropriate download/browse link       |
| `/browse`      | Directory listing (directory mode only)                                     |
| `/browse/`     | Subdirectory listing (directory mode only)                                  |
| `/download/`   | File download endpoint (works for both single file and directory modes)     |

## Technical Details

- **Port**: 8080 (hardcoded)
- **Token length**: 32 bytes (hex encoded)
- **Token validity**: 1 hour
- **Rate limit**: 100 requests/minute
- **Request size limit**: 10MB
- **Timeouts**:
  - Read/Write: 15 seconds
  - Idle: 30 seconds

## Building from Source

1. Ensure you have Go installed
2. Clone the repository
3. Build:
   ```sh
   go build -o server
   ```

## License

This project is open-source. Feel free to use and modify it according to your needs.

## Security Considerations

- Always run behind a reverse proxy in production
- Consider adding TLS/HTTPS
- The server binds to all interfaces by default - restrict as needed
- Tokens are stored in memory and cleared after expiration

## Known Limitations

- No persistent user accounts - tokens are temporary
- No upload functionality - read-only server
- Basic UI - functional but minimal styling

For any issues or feature requests, please open an issue on the repository.
