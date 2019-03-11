package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/common.git/config"
)

var rootCmd = &cobra.Command{
	Use:              "worker",
	TraverseChildren: true,
	Version:          "v0.1.0",
}

var cmdExample = &cobra.Command{
	Use:   "example [OPTIONS]",
	Short: "An example command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Log-Level:", viper.GetString("log.level"))
		fmt.Println("Log-Format:", viper.GetString("log.format"))
		fmt.Println("Eth-Clients:", viper.GetStringSlice("eth.clients"))
	},
}

func init() {
	rootCmd.AddCommand(cmdExample)
	config.LogLevel(cmdExample.Flags())
	config.LogFormat(cmdExample.Flags())
	config.EthClientURLs(cmdExample.Flags())
}

func main() {
	rootCmd.Execute()
}
