package main

import (
	"github.com/spf13/cobra"

	"github.com/syhily/hobbit/version"
)

var serverCmd = &cobra.Command{
	Use:   "hobbit-server",
	Short: "Hobbit is a proxy service to expose the Synology server to public networks.",
	Long: `A well designed proxy built with love by syhily and friends in Go.
                Complete documentation is available at https://syhily.github.io/hobbit`,
	Run: func(cmd *cobra.Command, args []string) {
		// This main function.
	},
}

func init() {
	serverCmd.AddCommand(version.VersionCmd)
}

func main() {
	if err := serverCmd.Execute(); err != nil {
		panic(err)
	}
}
