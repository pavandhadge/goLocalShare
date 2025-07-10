# Web UI & API Reference

---

## Web UI Overview
- **Owner Panel (`/`):** Shows access token, lets you reset it, and provides a link to the user page.
- **User Page (`/token`):** Enter the token to access shared files or directories.
- **Browse (`/browse`):** Directory listing (if sharing a directory).
- **Download (`/file` or `/api/download/`):** Download files (token required).

---

## API Endpoints

| Endpoint           | Method | Auth Required | Description                                 |
|--------------------|--------|--------------|---------------------------------------------|
| `/api/files`       | GET    | Yes          | List files/directories                      |
| `/api/download/`   | GET    | Yes          | Download a file                             |
| `/token`           | GET    | No           | User access page (HTML)                     |
| `/`                | GET    | No           | Owner panel (HTML)                          |
| `/reset-token`     | POST   | No           | Reset the access token (owner only)         |

---

## Example: List Files

**Request:**
```http
GET /api/files
X-Auth-Token: <token>
```

**Response:**
```json
{
  "files": [
    { "name": "file.txt", "isDirectory": false, "size": 1234, ... },
    ...
  ],
  "currentDir": "/",
  "basePath": "/home/user/share",
  "baseInfo": { ... }
}
```

---

## Example: Download File

**Request:**
```http
GET /api/download/file.txt
X-Auth-Token: <token>
```

**Response:**
- File download (binary data)
- 401 Unauthorized if token is missing or invalid

---

## UI Flows
- **Owner:** Start server → See token on `/` → Share link + token → Reset token if needed
- **User:** Visit link → Enter token on `/token` → Browse/download files

---

## Error Handling
- 401 Unauthorized: Invalid or expired token
- 404 Not Found: File or directory does not exist
- 429 Too Many Requests: Rate limit exceeded

---

For advanced configuration, see [Configuration & Advanced Options](./configuration.md). 