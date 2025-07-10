[📚 Full Documentation](./docs/README.md)

# goLocalShare

A simple, secure, and fast file sharing server for your local network or the cloud.

---

## 🚀 What is goLocalShare?

goLocalShare lets you share files or directories from your computer with anyone on your local network (or via a cloud link) in seconds. It’s designed for privacy, security, and ease of use—no accounts, no setup, just a single command.

---

## ✨ Features

- **Easy File & Directory Sharing:** Share any file or folder instantly with a single command.
- **Secure Token Authentication:** Only users with your unique, time-limited token can access your files.
- **Time-limited Sessions:** Access automatically expires after your chosen duration (default: 1 hour).
- **Cloud Upload & Auto-Delete:** Optionally upload files to Cloudinary for remote sharing. Files are deleted from the cloud after your session ends.
- **Strong Security:**
  - CSP, XSS, and other secure HTTP headers
  - Path and symlink protection
  - Rate limiting and brute-force protection
- **No Size Limit:** Share files of any size (limited only by your network or Cloudinary plan).
- **Cross-platform:** Works on Linux, Windows, and macOS.

---

## 🛠️ Installation

### 1. Download a Release
- Go to [Releases](https://github.com/pavandhadge/goLocalShare/releases) and download the binary for your OS.
- Or, use the build script to cross-compile for any platform:

```sh
./build.sh
```

### 2. Build from Source (requires Go 1.16+)

```sh
git clone https://github.com/pavandhadge/goLocalShare.git
cd goLocalShare
go build -o goLocalShare main.go
```

---

## ⚡ Usage

### Share a Single File for 2 Hours
```sh
./goLocalShare --duration 2h /path/to/your/file.ext
```

### Share a Directory for 30 Minutes
```sh
./goLocalShare --dir --duration 30m /path/to/your/directory
```

### Upload a File to Cloudinary for 1 Hour
- **First time only:** Add your Cloudinary credentials with `--cloud-name`, `--cloud-key`, and `--cloud-secret`.
- Credentials are saved in `~/.gofileserver_cloudinary.json` for future use.

```sh
./goLocalShare --cloud --cloud-name <name> --cloud-key <key> --cloud-secret <secret> --duration 1h /path/to/your/file.ext
```
- **Next time:** Just use `--cloud` (credentials are loaded automatically):
```sh
./goLocalShare --cloud --duration 1h /path/to/your/file.ext
```

---

## 🔒 Security Details

- **Token Authentication:**
  - Each session generates a unique access token.
  - Only users with the token can browse or download files.
  - Token is valid for the session duration (default: 1 hour).
- **Rate Limiting:** 100 requests/minute per IP.
- **Brute-force Protection:** 2-second delay after failed attempts.
- **Security Headers:** CSP, XSS, no-sniff, and more.
- **Path & Symlink Protection:** Only files within the shared path are accessible.
- **Session Expiry:** All access is revoked after the session ends.

---

## 🌐 How It Works

1. **Start the server:** Run the command with your file or directory.
2. **Get the link and token:** The server prints a link and a unique token.
3. **Share with others:** Give the link and token to anyone on your network.
4. **Access:** Users enter the token to browse/download files.
5. **Session ends:** After the duration, access is automatically revoked.

---

## ☁️ Cloud Upload (Cloudinary)

- Use `--cloud` to upload a file to Cloudinary for remote sharing.
- File is deleted from Cloudinary after the session duration.
- Credentials are stored in `~/.gofileserver_cloudinary.json` after first use.
- You can update credentials by re-running with the flags.

---

## 🖥️ API Endpoints (for advanced users)

| Endpoint         | Description                                 |
|------------------|---------------------------------------------|
| `/`              | Owner panel: shows your access token         |
| `/token`         | User page: enter token to access files       |
| `/api/files`     | List files/directories (requires token)      |
| `/api/download/` | Download a file (requires token)             |

---

## 🧩 Troubleshooting

- **Port in use?** The server uses port 8090 by default. Stop other services or change the port in `main.go` if needed.
- **Cloudinary errors?** Double-check your credentials. Delete `~/.gofileserver_cloudinary.json` to reset.
- **Token not working?** Make sure you’re using the latest token and the session hasn’t expired.
- **Firewall issues?** Ensure port 8090 is open on your network.

---

## 📝 License

MIT. Free for personal and commercial use. Contributions welcome!

---

## 🙋 FAQ

**Q: Can I share with people outside my local network?**
- Yes, if your network/firewall allows, or by using the Cloudinary upload feature.

**Q: Are my files ever stored on a third-party server?**
- Only if you use `--cloud`. Otherwise, files stay on your machine.

**Q: How do I revoke access immediately?**
- Use the "Reset Token" button on the owner panel (`/`).

**Q: Is there a web UI?**
- Yes, a minimal web UI is provided for both owners and users.

---

## 👨‍💻 Author

Made with ❤️ by [Pavan Dhadge](https://github.com/pavandhadge)

For issues or suggestions, open an [issue](https://github.com/pavandhadge/goLocalShare/issues) or [pull request](https://github.com/pavandhadge/goLocalShare/pulls).
