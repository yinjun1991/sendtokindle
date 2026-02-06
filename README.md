# sendtokindle

Upload eBooks from your computer and download them on your Kindle (or any device on the same Wi‑Fi).

## Quick Start

1. Start the server:

   ```bash
   go run ./cmd/sendtokindle --port 8090
   ```

2. Open the admin page on your computer:

   - http://localhost:8090/admin

3. On the admin page:
   - Upload your `.epub` / `.pdf` files
   - Copy the **Kindle** link shown there and open it in the Kindle browser
   - Tap **Download** on the book you want

## Options

- `--port`  
  HTTP port to listen on (default: `8080`).

- `--dir`  
  Folder to store uploaded books (default: `~/.sendtokindle`).

Examples:

```bash
# Use a custom storage folder
go run ./cmd/sendtokindle --port 8090 --dir "$HOME/.sendtokindle"

# Use another port
go run ./cmd/sendtokindle --port 8081
```

