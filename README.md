# Nudgen CLI

`nudgen-cli` is the command-line command center for automating campaign email marketing. Effortlessly launch campaigns, manage audience segments, and scale your outreach through a developer-first interface built for both humans and AI agents.

- **Campaign Orchestration**: Deploy, monitor, and scale email marketing campaigns with terminal efficiency.
- **Agent Native**: Built from the ground up for seamless execution by AI agents like Claude, Codex, and Gemini.
- **Flexible CRM**: Import contact lists and manage brand identities across multi-tenant team environments.
- **Real-Time Insights**: Monitor campaign activity, referral conversions, and payout history directly.

## Quick Start

```bash
curl -fsSL https://raw.githubusercontent.com/Nudgen-Marketing/nudgen-cli/main/scripts/install.sh | bash
```

That's it. You now have full access to Nudgen from your terminal.

<details>
<summary>Other installation methods</summary>

**Go install:**
```bash
go install github.com/Nudgen-Marketing/nudgen-cli@latest
```

**Local build:**
```bash
make build
./bin/nudgen --help
```

**GitHub Release:** Download binaries for macOS, Linux, and Windows from [Releases](https://github.com/Nudgen-Marketing/nudgen-cli/releases).

</details>

## Usage

```bash
nudgen login                  # Authenticate via Personal Access Token
nudgen whoami                 # Check current identity
nudgen teams list             # List your teams
nudgen teams switch <id>      # Switch to a different team context
nudgen brand list             # List brand identities for the team
nudgen contacts list          # List active contacts
nudgen contacts create        # Create a new contact (interactive)
nudgen contacts update <id>   # Update an existing contact (interactive)
nudgen contacts delete <id>   # Delete a contact
nudgen contacts import        # Import contacts from CSV (interactive)
nudgen campaigns list         # List active campaigns in current team
nudgen referral activity      # View recent referral click history
nudgen referral payouts       # View your payout history
```

### JSON Output
Every command supports the `--json` flag for machine-readable output:

```bash
nudgen campaigns list --json
```

## Configuration

The CLI stores its state in your home directory:

```
~/.nudgen/
├── config.json               # Persistent team context (activeTeamId)
└── (system keychain)         # Secure storage for your PAT token
```

## Development

```bash
make build            # Build binary to ./bin/nudgen
make test             # Run Go tests
make fmt              # Format code
make vet              # Run go vet
make help             # Show all Makefile targets
```

See [AGENTS.md](AGENTS.md) for internal project structure and development context.

## License

[MIT](MIT-LICENSE)
