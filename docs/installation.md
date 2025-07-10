# Installation & Building

---

## Downloading a Release
- Visit [Releases](https://github.com/pavandhadge/goLocalShare/releases) and download the binary for your OS (Linux, Windows, macOS).
- Extract the archive if needed and place the binary somewhere in your PATH.

---

## Building from Source

### Prerequisites
- Go 1.16 or newer
- Git (for cloning the repo)

### Steps
```sh
git clone https://github.com/pavandhadge/goLocalShare.git
cd goLocalShare
go build -o goLocalShare main.go
```
- The binary `goLocalShare` will be created in the current directory.

---

## Cross-Compiling (All Platforms)

A build script is provided to build for all major platforms:
```sh
./build.sh
```
- Output binaries are placed in the `builds/` directory.
- Checksums are generated for verification.

---

## Verifying the Build
- Run `./goLocalShare --help` to check that the binary works.
- If you see usage instructions, the build was successful.

---

## Troubleshooting
- **Go not found?** Install Go from [golang.org/dl](https://golang.org/dl/).
- **Permission denied?** Run `chmod +x goLocalShare` to make the binary executable.
- **Build errors?** Ensure you have Go 1.16+ and your `GOPATH`/`GOROOT` are set correctly.
- **Windows:** Use `goLocalShare.exe` instead of `./goLocalShare`. 