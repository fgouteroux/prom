package version

// Default build-time variable.
// These values are overridden via ldflags
var (
	ProgName     = "prom"
	PlatformName = "unknown-platform"
	Version      = "unknown-version"
	GitCommit    = "unknown-commit"
	BuildTime    = "unknown-buildtime"
)
