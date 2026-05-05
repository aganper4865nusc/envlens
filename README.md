# envlens

> Utility to diff, validate, and audit environment variable files across deployment stages.

---

## Installation

```bash
go install github.com/yourname/envlens@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/envlens.git && cd envlens && go build ./...
```

---

## Usage

**Diff two environment files:**

```bash
envlens diff .env.staging .env.production
```

**Validate a file against a required keys schema:**

```bash
envlens validate .env --schema .env.schema
```

**Audit for missing, extra, or mismatched variables across stages:**

```bash
envlens audit --stages .env.dev,.env.staging,.env.production
```

Example output:

```
[MISSING]  .env.production  → DB_TIMEOUT
[EXTRA]    .env.staging     → DEBUG_MODE
[MISMATCH] LOG_LEVEL        → dev=debug | staging=info | production=info
```

---

## Flags

| Flag | Description |
|------|-------------|
| `--schema` | Path to a schema file defining required keys |
| `--stages` | Comma-separated list of env files to audit |
| `--quiet` | Suppress output, exit code only |
| `--json` | Output results as JSON |

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE) © yourname