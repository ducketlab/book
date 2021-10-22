package cmd

import (
	"errors"
	"fmt"
	"github.com/ducketlab/book/version"
	"github.com/spf13/cobra"
	"os"
)

var vers bool

var RootCmd = &cobra.Command{
	Use: "book",
	RunE: func(cmd *cobra.Command, args []string) error {
		if vers {
			fmt.Println(version.FullVersion())
			return nil
		}
		return errors.New("no flags find")
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&vers,
		"version", "v", false, "the version")
}
