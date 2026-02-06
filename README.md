# sendtokindle

Upload eBooks from your computer and download them on your Kindle (or any device on the same Wi‑Fi).

## Quick Start

1. Download the binary for your OS from the GitHub Releases page.

2. Start the server:

   ```bash
   ./sendtokindle --port 8090
   ```

3. Open the admin page on your computer:

   - http://localhost:8090/admin

4. On the admin page:
   - Upload files your Kindle can open (for example: `.azw3`, `.azw`, `.mobi`, `.txt`)
   - Copy the **Kindle** link shown there and open it in the Kindle browser
   - Tap **Download** on the book you want

## Options

- `--port`  
  Port to listen on (default: `8080`).

- `--dir`  
  Folder to store uploaded books (default: `~/.sendtokindle`).

## macOS: first run

If macOS blocks the binary with a security warning, you can allow it once:

- Finder: right click the `sendtokindle` binary → **Open** → **Open**
- Or remove the quarantine attribute:

```bash
xattr -d com.apple.quarantine ./sendtokindle
```

Examples:

```bash
# Use a custom storage folder
 ./sendtokindle --port 8090 --dir "$HOME/.sendtokindle"

# Use another port
 ./sendtokindle --port 8081
```
