# HashGen — Hash Generating Service
In this week, I built a 10-character unique hash generation service that accepts alphanumeric input. The solution uses SHA-256 for secure hashing and Base62 encoding to produce compact, URL-safe output. The generated hash can also be easily copied for reuse.

## Design Decisions

### Why SHA-256?

I used SHA-256 as the hash function for the following reasons:

- **Collision resistance**: SHA-256 is a one-way function with low probability of two different inputs producing the same hash, making it reliable for unique identifier generation.
- **Determinism**: The same input always yields the same 256-bit digest, which is required for a hash service to be useful and testable.
- **Avalanche effect**: A single character change in the input flips roughly half of the output bits, ensuring even similar inputs produce very different hashes. This property is validated in the test suite.
- **Industry standard**: SHA-256 is widely used, part of the SHA-2 family standardized by NIST, and implemented natively in Go's `crypto/sha256` package — no external dependencies needed.
- **Speed**: It is computationally fast for short inputs while remaining cryptographically strong.

---

### Why Base62 Encoding? 

> The encoding I have used here is **Base62**.

Base62 uses the alphabet `0–9`, `A–Z`, `a–z` (62 characters total). Base64 introduces `+`, `/`, and `=` characters which are URL-unsafe and require escaping in query strings, paths, and JSON.

Since the output is intended to be a short, shareable, URL-safe identifier (10 characters), Base62 is the more practical choice. The first 8 bytes of the SHA-256 digest are interpreted as a `uint64` and then Base62-encoded, giving a compact alphanumeric string that is both URL-safe and human-readable.

---

### Approach to Reach the Output

The hash generation follows this pipeline:

```
User Input (alphanumeric string)
        │
        ▼
  Input Validation
  (only a-z, A-Z, 0-9 accepted)
        │
        ▼
  SHA-256 Hash
  (produces 32-byte / 256-bit digest)
        │
        ▼
  Extract first 8 bytes
  (reduces 32 bytes → 8 bytes)
        │
        ▼
  Convert to uint64
  (big-endian binary interpretation)
        │
        ▼
  Base62 Encode
  (converts uint64 to alphanumeric string)
        │
        ▼
  Pad or Truncate to 10 characters
  (ensures fixed-length output)
        │
        ▼
  Final Hash (10-character string)
```

**Why only 8 bytes of SHA-256?**
The full 32-byte SHA-256 digest would produce a very long Base62 string (~43 chars). For a compact, fixed-length identifier, only the first 8 bytes (64 bits) are used. This still provides 62^10 ≈ 839 trillion unique values — more than sufficient for typical use cases, while keeping the hash short and readable.

---

## Project Structure

```
week2/
└── hash-service/
    ├── main.go                  # Entry point, HTTP server setup, graceful shutdown
    ├── go.mod                   # Go module (no external dependencies)
    ├── config/
    │   └── config.go            # Port configuration via environment variable
    ├── handlers/
    │   └── handlers.go          # HTTP route handlers and middleware
    ├── hashgen/
    │   ├── generator.go         # Core hash generation logic
    │   ├── base62.go            # Base62 encoding utility
    │   └── generator_test.go    # Unit tests (determinism, length, avalanche effect)
    ├── templates/
    │   └── index.html           # Web UI (embedded into binary)
    └── static/
        └── style.css            # Stylesheet (embedded into binary)
```

---

**Validation rules**
- Input must be non-empty
- Input must contain only alphanumeric characters (`a-z`, `A-Z`, `0-9`)

---

## Getting Started

### Prerequisites

- Go 1.22+

### Run locally

```bash
cd hash-service
go run main.go
```

The server starts on port `8080` by default.

```bash
# Use a custom port
PORT=9000 go run main.go
```

### Build a binary

```bash
go build -o hashgen
./hashgen
```

Static assets and HTML templates are embedded directly into the binary — no extra files needed at runtime.

### Run tests

```bash
go test ./hashgen/...
```

Tests cover:
- **Determinism** — same input always produces the same hash
- **Fixed length** — output is always exactly 10 characters
- **Avalanche effect** — small input changes produce very different hashes

---

## Example Usage

**Web UI**: open `http://localhost:8080` in your browser.

---

## Key Properties

| Property | Detail |
|---|---|
| Hash algorithm | SHA-256 |
| Output encoding | Base62 (`0-9A-Za-z`) |
| Output length | Fixed 10 characters |
| Input constraints | Alphanumeric only |
| External dependencies | None (stdlib only) |
| Graceful shutdown | Yes (SIGINT / SIGTERM) |
| Static asset delivery | Embedded in binary |