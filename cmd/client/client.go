package main

import (
	"github.com/spf13/cobra"

	"github.com/syhily/hobbit/version"
)

var clientCmd = &cobra.Command{
	Use:   "hobbit-client",
	Short: "Client to connect with the Hobbit server. Installed on Synology.",
	Long: `A well designed proxy built with love by syhily and friends in Go.
                Complete documentation is available at https://syhily.github.io/hobbit`,
	Run: func(cmd *cobra.Command, args []string) {
		// This main function.
	},
}

func init() {
	clientCmd.AddCommand(version.VersionCmd)
}

func main() {
	if err := clientCmd.Execute(); err != nil {
		panic(err)
	}
}
