# Configuration & Advanced Options

---

## Changing the Port
- By default, goLocalShare uses port 8090.
- To change the port, edit the `port` variable in `main.go` and rebuild:
  ```go
  port := ":8090" // Change to your desired port
  ```

---

## Environment Variables
- goLocalShare does not use environment variables by default, but you can set them in your shell for custom scripts or wrappers.

---

## Config Files
- Cloudinary credentials are stored in `~/.gofileserver_cloudinary.json` after first use of cloud mode.
- To reset, delete this file and re-run with the `--cloud-*` flags.

---

## Custom Builds
- You can modify the source code to change defaults, add features, or customize the UI.
- Use the provided `build.sh` script for cross-compilation.

---

## Advanced Flags
- All CLI flags can be combined as needed.
- For example, share a directory in cloud mode for 2 hours:
  ```sh
  ./goLocalShare --dir --cloud --duration 2h /path/to/dir
  ```

---

For troubleshooting, see [Troubleshooting & FAQ](./troubleshooting.md). 