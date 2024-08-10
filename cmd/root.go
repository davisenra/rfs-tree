package cmd

import (
	"fmt"
	"os"

	"github.com/davisenra/rfs-tree/internal/tree"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rfs-tree <path>",
	Short: "rfs-tree is a cli tool for displaying file system trees",
	Long:  "rfs-tree is a cli tool for displaying file system trees",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		root := args[0]
		node, err := tree.GenerateTree(root)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		tree.OutputTree(node, os.Stdout)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing Zero '%s'\n", err)
		os.Exit(1)
	}
}
