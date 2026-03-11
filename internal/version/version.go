package version

import (
	"fmt"
	"runtime"
)

var (
	// Version is the current version (set via ldflags)
	Version = "dev"
	// GitCommit is the git commit hash (set via ldflags)
	GitCommit = "none"
	// BuildDate is the build date (set via ldflags)
	BuildDate = "unknown"
)

// Info contains version information
type Info struct {
	Version   string
	GitCommit string
	BuildDate string
	GoVersion string
	Platform  string
}

// Get returns the current version info
func Get() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String returns a formatted version string
func (i Info) String() string {
	return fmt.Sprintf("gps version %s\n  commit: %s\n  built: %s\n  go: %s\n  platform: %s",
		i.Version, i.GitCommit, i.BuildDate, i.GoVersion, i.Platform)
}
