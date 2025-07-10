# Quick Start Guide

Get up and running with goLocalShare in minutes!

---

## Prerequisites
- Go 1.16+ (for building from source)
- Linux, Windows, or macOS

---

## 1. Download or Build

**Option A: Download a Release**
- Go to [Releases](https://github.com/pavandhadge/goLocalShare/releases) and download the binary for your OS.

**Option B: Build from Source**
```sh
git clone https://github.com/pavandhadge/goLocalShare.git
cd goLocalShare
go build -o goLocalShare main.go
```

---

## 2. Share Your First File

```sh
./goLocalShare --duration 1h /path/to/your/file.ext
```

- The server will print a link and a unique access token.
- Open the link in your browser and enter the token to access the file.

---

## 3. Share a Directory

```sh
./goLocalShare --dir --duration 30m /path/to/your/directory
```

---

## 4. (Optional) Cloud Upload

```sh
./goLocalShare --cloud --cloud-name <name> --cloud-key <key> --cloud-secret <secret> --duration 1h /path/to/your/file.ext
```
- Credentials are saved for future use.

---

## 5. Stop Sharing
- The session ends automatically after the duration.
- Or, press `Ctrl+C` in the terminal to stop the server immediately.

---

For more, see [Usage & CLI Reference](./usage.md). 