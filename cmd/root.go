/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"http-spy-forward/modules"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "http-spy-forward",
	Short: "\nA tool for sniff HTTP packets and forward to the target server.",
	Run:   modules.Start,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("iface", "i", "", "Network interface to sniff")
	rootCmd.Flags().BoolP("debug", "d", false, "Debug mode")
	rootCmd.Flags().StringP("filter", "f", "", "Filter expression to sniff")
	rootCmd.Flags().Int32P("length", "l", 1024, "Length of packet to sniff")
	rootCmd.Flags().StringP("url", "u", "", "URL to forward request to")

	rootCmd.MarkFlagRequired("iface")
	rootCmd.MarkFlagRequired("url")
}
