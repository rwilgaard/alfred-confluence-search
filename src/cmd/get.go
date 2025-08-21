package cmd

import (
	"github.com/spf13/cobra"
)

var (
    getCmd = &cobra.Command{
        Use:   "get",
        Short: "get one or many resources",
        GroupID: groupBase,
    }
)

func init() {
    getCmd.AddGroup(
        &cobra.Group{
            ID: groupPages,
            Title: "Page commands:",
        },
        &cobra.Group{
            ID: groupSpaces,
            Title: "Space commands:",
        },
    )
    rootCmd.AddCommand(getCmd)
}
