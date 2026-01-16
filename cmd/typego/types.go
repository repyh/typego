package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var typesCmd = &cobra.Command{
	Use:   "types",
	Short: "Sync and update TypeGo ambient definitions",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Syncing TypeGo definitions...")

		dtsPath := filepath.Join(".typego", "types", "go.d.ts")
		if err := os.MkdirAll(filepath.Dir(dtsPath), 0755); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		masterDts, err := os.ReadFile("go.d.ts")
		if err != nil {
			fmt.Println("Warning: master go.d.ts not found in root.")
			return
		}

		if err := os.WriteFile(dtsPath, masterDts, 0644); err != nil {
			fmt.Printf("Error writing types: %v\n", err)
			return
		}

		fmt.Println("✅ Definitions synced to .typego/types/go.d.ts")
	},
}

func init() {
	rootCmd.AddCommand(typesCmd)
}
