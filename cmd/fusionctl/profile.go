package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/natesales/openreactor/pkg/profile"
)

var (
	file string
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage profiles",
}

var profileGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new profile with default values",
	Run: func(cmd *cobra.Command, args []string) {
		defaultProfile, err := profile.Parse([]byte("name: Default\nrevision: 0"))
		if err != nil {
			fmt.Println(err)
			return
		}

		y, err := yaml.Marshal(defaultProfile)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(y))
	},
}

var profileLintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint a profile",
	Run: func(cmd *cobra.Command, args []string) {
		contents, err := os.ReadFile(file)
		if err != nil {
			fmt.Println(err)
			return
		}

		p, err := profile.Parse(contents)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%s (rev. %s) is valid\n", p.Name, p.Revision)
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)

	profileCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "Profile file")
	if err := profileCmd.MarkPersistentFlagRequired("file"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	profileCmd.AddCommand(profileGenerateCmd)
	profileCmd.AddCommand(profileLintCmd)
}
