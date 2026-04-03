# Nudgen CLI Agent Skill

Master the Nudgen marketing automation platform via its interactive CLI. This skill allows managing campaigns, contacts, stats, and referrals.

## Tools

Run the following commands via standard terminal execution.

### Authentication & Identification
- `nudgen login`: Authenticate with a PAT (Interactive).
- `nudgen whoami`: Show current token owner and user info.
- `nudgen logout`: Clear credentials.

### Team Context Management
- `nudgen teams list`: List all teams the user belongs to.
- `nudgen teams current`: Show currently active team.
- `nudgen teams switch <id>`: Switch context to a specific team.

### Campaigns
Manage marketing outreach efforts within the active team.
- `nudgen campaigns list`: List active campaigns.
- `nudgen campaigns create`: Interactive creation wizard.
- `nudgen campaigns update <id>`: Modify existing campaign settings.
- `nudgen campaigns delete <id>`: Remove a campaign record.

### Lead Management (Contacts)
Manage your team's lead database and segmentation.
- `nudgen contacts list`: List all leads in the current team.
- `nudgen contacts create`: Manually add a new lead record.
- `nudgen contacts import`: Interactive CSV mapper for bulk ingestion.
- `nudgen contacts update <id>`: Edit lead details.
- `nudgen contacts delete <id>`: Remove a contact record.

### Brand Identities
Control AI Agents' representation of your company.
- `nudgen brand list`: List all brand voices for the team.
- `nudgen brand create`: Wizard to initialize a new brand voice.
- `nudgen brand update <id>`: Edit brand configurations.
- `nudgen brand delete <id>`: Remove a brand identity.

### Referrals (Earn)
- `nudgen referral stats`: View high-level referral performance summary.
- `nudgen referral link`: Get your unique referral link.
- `nudgen referral activity`: View recent click history and signups.
- `nudgen referral referrals`: List teams referred by you and their status.
- `nudgen referral payouts`: View history of payout requests and status.

## Interaction Patterns

### JSON Automation
When automating as an agent, **always** append `--json` to commands.
- `nudgen campaigns list --json`
- `nudgen teams list --json`

### Multi-tenant Workflow
Before performing any data operations (campaigns, brands), ensure you check the active team context:
1. Run `nudgen teams current --json`. 
2. If it's the wrong team, run `nudgen teams switch <target-team-id>`.
3. Proceed with data commands.

## Architecture
The CLI is written in Go and communicates with the REST API at `nudgen.cc/api/v1`.
Credentials are encrypted in the system keychain, while non-sensitive state like the `active_team_id` is in `~/.nudgen/config.json`.
