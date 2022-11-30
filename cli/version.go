package cli

import (
	"fmt"
	"runtime"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/uor-framework/uor-fuse-go/config"
)

var (
	// commit is the head commit from git
	commit string
	// buildDate in ISO8601 format
	buildDate string
	// version describes the version of the client
	// set at build time or detected during runtime.
	version = "v0.0.0-unknown"
	// buildData set at build time to add extra information
	// to the version.
	buildData string
)

var versionTemplate = `UOR FUSE Driver:
 Version:	{{ .Version }}
 Go Version:	{{ .GoVersion }}
 Git Commit:	{{ .GitCommit }}
 Build Date:	{{ .BuildDate }}
 Platform:	{{ .Platform }}
`

type clientVersion struct {
	Platform  string
	Version   string
	GitCommit string
	GoVersion string
	BuildDate string
}

// NewVersionCmd creates a new cobra.Command for the version subcommand.
func NewVersionCmd(rootOpts *config.RootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return getVersion(rootOpts)
		},
	}
}

// getVersion will output the templated version message.
func getVersion(ro *config.RootOptions) error {
	versionWithBuild := func() string {
		if buildData != "" {
			return fmt.Sprintf("%s+%s", version, buildData)
		}
		return version
	}

	versionInfo := clientVersion{
		Version:   versionWithBuild(),
		GitCommit: commit,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}

	tmp, err := template.New("version").Parse(versionTemplate)
	if err != nil {
		return fmt.Errorf("template parsing error: %v", err)
	}

	return tmp.Execute(ro.IOStreams.Out, versionInfo)
}
