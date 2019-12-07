// Package cmd provides the command line interface to scrape cannaconnection.com
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "strains",
	Short: "Scrapes strain data from https://www.cannaconnection.com/strains",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Set author information.
	rootCmd.PersistentFlags().StringP("author", "a", "Dani Roxberry (@droxey)", "Author name for copyright attribution")
}
