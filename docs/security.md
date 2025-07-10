# Security Model

---

## Authentication & Token System
- Each sharing session generates a unique, random access token (32 bytes, hex encoded).
- All API and file access requires this token, sent via header or query parameter.
- Tokens are valid only for the configured session duration (default: 1 hour).
- Tokens can be reset by the owner at any time, immediately revoking access.

---

## Session Expiry
- When the session duration ends, all tokens are invalidated and access is denied.
- The server enforces session expiry for all endpoints.

---

## Rate Limiting & Brute-Force Protection
- 100 requests per minute per IP address.
- After a failed authentication attempt, a 2-second delay is enforced for that IP.
- These measures help prevent brute-force and denial-of-service attacks.

---

## Security Headers
- Content-Security-Policy (CSP): Restricts sources for scripts and styles.
- X-Content-Type-Options: Prevents MIME sniffing.
- X-Frame-Options: DENY (prevents clickjacking).
- X-XSS-Protection: Enables browser XSS filter.
- Referrer-Policy: no-referrer.

---

## Path & Symlink Protection
- All file and directory access is restricted to the shared base path.
- Symlinks are checked to ensure they do not escape the base directory.
- Directory traversal attacks are blocked by path sanitization.

---

## Threat Model
- **Attacker on LAN:** Cannot access files without the token.
- **Token leakage:** Anyone with the token can access files until session expiry or token reset.
- **Cloud mode:** Files are uploaded to Cloudinary and deleted after the session. Cloudinary credentials are stored locally and never sent to others.

---

## Deployment Best Practices
- Run behind a reverse proxy (e.g., Nginx) for additional security and TLS/HTTPS.
- Restrict server to listen only on trusted interfaces if not sharing with the whole LAN.
- Use strong, unique Cloudinary credentials for cloud mode.
- Regularly update goLocalShare to receive security patches.

---

## Recommendations
- Never share your access token publicly.
- Reset the token immediately if you suspect it has leaked.
- Use the shortest session duration necessary for your use case.
- For sensitive files, prefer LAN mode over cloud mode. 