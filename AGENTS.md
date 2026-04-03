# AGENTS.md

## Repository Structure

```
nudgen-cli/
├── cmd/               # Command-line interface definitions (Cobra)
├── internal/
│   ├── api/           # API Client and data models
│   ├── config/        # Config and local team-state management
│   └── version/       # Build metadata (Version, Commit, Date)
├── scripts/           # Premium shell installer tools
├── docs/              # Project design plans and architecture
└── Makefile           # High-level automation targets
```

## Internal Architecture

The CLI follows the "Interaction over API" principle. It uses standard RESTful API calls to communicate with the Nudgen backend, while securely managing session state in the user's localized configuration.

### Important Patterns

1.  **Context-Aware**: The `active_team_id` is stored in `~/.nudgen/config.json`. Any command fetching data (campaigns, brands, etc.) MUST assume it is relative to the active team currently stored there.
2.  **Secure Storage**: PAT (Personal Access Token) is stored in the system keychain via `github.com/zalando/go-keyring`. Never log or output the raw token.
3.  **JSON Standard**: All commands must support `--json`. This is critical for agent-to-agent communication.

## Development Context

### Build Pipeline

Nudgen CLI uses `goreleaser` for automated releases on GitHub tags (`v*`). The binary is stamped with version data at build time using `-X` ldflags in the `Makefile`.

### Quality Gate

Before pushing changes, run:

```bash
make fmt && make vet && make test
```

### Adding New Commands

1.  Add necessary models to `internal/api/models.go`.
2.  Add API methods to `internal/api/client.go`.
3.  Create the command in `cmd/`.
4.  Register the command in `cmd/root.go`'s `init()` if it's a top-level command.
