package version

import "fmt"

// Version is the current version of the CLI, set at build time via ldflags.
var Version = "dev"

// Commit is the git commit hash, set at build time.
var Commit = "none"

// Date is the build date, set at build time.
var Date = "unknown"

func FullVersion() string {
	return fmt.Sprintf("nudgen %s (commit: %s, date: %s)", Version, Commit, Date)
}
