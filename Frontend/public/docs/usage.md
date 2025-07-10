# Usage & CLI Reference

---

## Basic Usage

Share a file for 1 hour:
```sh
./goLocalShare --duration 1h /path/to/your/file.ext
```

Share a directory for 30 minutes:
```sh
./goLocalShare --dir --duration 30m /path/to/your/directory
```

Upload a file to Cloudinary for 1 hour:
```sh
./goLocalShare --cloud --cloud-name <name> --cloud-key <key> --cloud-secret <secret> --duration 1h /path/to/your/file.ext
```

---

## CLI Flags & Options

| Flag                | Description                                              | Example                                  |
|---------------------|----------------------------------------------------------|------------------------------------------|
| `--duration <time>` | Session duration (e.g. 1h, 30m, 2h30m)                  | `--duration 2h`                          |
| `--dir`             | Share a directory instead of a single file               | `--dir`                                  |
| `--cloud`           | Upload file to Cloudinary for remote sharing             | `--cloud`                                |
| `--cloud-name`      | Cloudinary cloud name (required for first cloud upload)  | `--cloud-name demo`                      |
| `--cloud-key`       | Cloudinary API key (required for first cloud upload)     | `--cloud-key 12345`                      |
| `--cloud-secret`    | Cloudinary API secret (required for first cloud upload)  | `--cloud-secret abcde`                   |

---

## Token System
- Each session generates a unique access token.
- Only users with the token can access files.
- Token is valid for the session duration.
- Token can be reset from the owner panel (`/`).

---

## Web UI Basics
- **Owner Panel (`/`):** Shows your access token, lets you reset it, and provides a link to the user page.
- **User Page (`/token`):** Enter the token to access shared files or directories.
- **Browse (`/browse`):** Directory listing (if sharing a directory).
- **Download (`/file` or `/api/download/`):** Download files (token required).

---

## Example Scenarios

| Scenario                        | Command Example                                                      |
|---------------------------------|----------------------------------------------------------------------|
| Share a file for 2 hours        | `./goLocalShare --duration 2h /path/to/file.ext`                     |
| Share a directory for 30 mins   | `./goLocalShare --dir --duration 30m /path/to/dir`                   |
| Cloud upload for 1 hour         | `./goLocalShare --cloud --cloud-name n --cloud-key k --cloud-secret s --duration 1h /path/to/file.ext` |
| Reset token                     | Use the "Reset Token" button on the owner panel                      |

---

For advanced configuration, see [Configuration & Advanced Options](./configuration.md). 