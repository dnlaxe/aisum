# aisum

My terminal gets messy. I want to make dumps more understandable by using AI to summarise piped output.

### Prerequisites

- Go 1.18 or higher.
- Ollama running locally.

### Build

1. Clone the repo:

   ```
   git clone https://github.com/dnlaxe/aisum.git
   cd aisum
   ```

2. Build + install the binary:

   ```
   make install
   ```

   This compiles the tool and moves it to `/usr/local/bin`, allowing you to run `aisum` from any directory.

### Use

```bash
    command | aisum
```

Pipe command output into aisum to get a shorter AI-generated summary.

### Examples

```bash
    go test ./... 2>&1 | aisum
    docker logs api | aisum
    git diff | aisum
```

### License

MIT
