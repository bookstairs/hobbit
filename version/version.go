package version

import (
	"github.com/spf13/cobra"

	"github.com/syhily/hobbit/config"
	"github.com/syhily/hobbit/pkg/logger"
)

// VersionCmd represents the version command.
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return the hobbit version info",
	Run: func(cmd *cobra.Command, args []string) {
		logger.NewTableLogger().
			Title("hobbit version info").
			Row("Version", config.GitVersion).
			Row("Commit", config.GitCommit).
			Row("Build Date", config.BuildTime).
			Row("Go Version", config.GoVersion).
			Row("Platform", config.Platform).
			Log()
	},
}
