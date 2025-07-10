# Cloud Upload & Integrations

---

## How Cloud Mode Works
- When `--cloud` is used, the file is uploaded to Cloudinary using your credentials.
- The server prints a Cloudinary URL for sharing.
- After the session duration, the file is automatically deleted from Cloudinary.

---

## Credential Management
- On first use, provide `--cloud-name`, `--cloud-key`, and `--cloud-secret`.
- Credentials are saved in `~/.gofileserver_cloudinary.json` for future use.
- To update credentials, re-run with the flags; the config file will be overwritten.
- Credentials are never sent to clients or included in shared links.

---

## File Lifecycle
1. File is uploaded to Cloudinary at session start.
2. A public (but unguessable) URL is generated.
3. After the session duration, the file is deleted from Cloudinary via API.

---

## Security & Privacy
- Only the owner’s machine and Cloudinary see the file.
- The Cloudinary URL is not guessable, but anyone with the link can access the file until deletion.
- Credentials are stored locally and never shared.

---

## Cloudinary Details
- Uses the [cloudinary-go](https://github.com/cloudinary/cloudinary-go) SDK.
- Supports all file types and sizes allowed by your Cloudinary plan.
- PublicID is randomized for each upload.

---

## Future Integrations
- The architecture allows for additional cloud providers (e.g., S3, Dropbox) in the future.
- Contributions for new integrations are welcome! 