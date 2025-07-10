# Introduction & Overview

## What is goLocalShare?

goLocalShare is a modern, open-source file sharing server designed for simplicity, security, and speed. It allows you to share files or directories from your computer with anyone on your local network—or, optionally, via a secure cloud link—using a single command. No accounts, no complex setup, just instant sharing.

## Philosophy
- **Simplicity:** One command to share, one link to access. Minimal configuration.
- **Security:** Every session is protected by a unique, time-limited token. Strong security headers, path/symlink protection, and rate limiting are built-in.
- **Privacy:** Files are never uploaded to third-party servers unless you explicitly use cloud mode.
- **Transparency:** Open source, auditable, and easy to understand.

## Who is it for?
- Home users who want to share files with family or friends on the same Wi-Fi.
- Developers and IT professionals who need a quick, secure way to transfer files between devices.
- Anyone who wants a portable, zero-setup file server for LAN or cloud sharing.

## Key Features
- Share any file or directory instantly over your local network.
- Optional cloud upload (Cloudinary) for remote sharing, with auto-delete after session.
- Secure, time-limited access tokens for every session.
- Minimal, responsive web UI for both owners and users.
- Strong security: CSP, XSS, rate limiting, brute-force protection, path/symlink checks.
- Cross-platform: Linux, Windows, macOS.
- No persistent user accounts or tracking.

## Project History
- **2023:** Initial version released as a simple Go file server for local sharing.
- **2024:** Major refactor for security, cloud upload, and modern web UI. Renamed to goLocalShare.
- **Ongoing:** Actively maintained with a focus on usability, security, and community feedback. 