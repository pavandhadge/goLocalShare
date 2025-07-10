# Troubleshooting & FAQ

---

## Common Issues & Solutions

### Port Already in Use
- **Error:** "address already in use"
- **Solution:** Stop other services using port 8090, or change the port in `main.go` and rebuild.

### Permission Denied
- **Error:** "permission denied"
- **Solution:** Run `chmod +x goLocalShare` to make the binary executable.

### Cloudinary Upload Fails
- **Error:** "Cloudinary credentials required" or upload fails
- **Solution:** Double-check your credentials. Delete `~/.gofileserver_cloudinary.json` to reset.

### Token Not Working
- **Error:** "Invalid or expired token"
- **Solution:** Make sure you’re using the latest token and the session hasn’t expired. Reset the token if needed.

### Firewall/Network Issues
- **Error:** Cannot access server from other devices
- **Solution:** Ensure port 8090 is open on your network and firewall.

---

## Advanced Debugging
- Run with `go run main.go ...` for live debugging.
- Check server logs for detailed error messages.
- Use browser dev tools to inspect API requests and responses.

---

## FAQ

**Q: Can I share with people outside my local network?**
- Yes, if your network/firewall allows, or by using the Cloudinary upload feature.

**Q: Are my files ever stored on a third-party server?**
- Only if you use `--cloud`. Otherwise, files stay on your machine.

**Q: How do I revoke access immediately?**
- Use the "Reset Token" button on the owner panel (`/`).

**Q: Is there a web UI?**
- Yes, a minimal web UI is provided for both owners and users.

**Q: How do I get help or report a bug?**
- Open an [issue](https://github.com/pavandhadge/goLocalShare/issues) on GitHub.

---

For more, see the [Documentation Index](./README.md). 