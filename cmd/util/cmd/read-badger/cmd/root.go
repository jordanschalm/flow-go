package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagDatadir   string
	flagPebbleDir string
)

var rootCmd = &cobra.Command{
	Use:   "read-badger",
	Short: "read storage data",
}

var RootCmd = rootCmd

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&flagDatadir, "datadir", "d", "/var/flow/data/protocol", "directory to the badger dababase")
	_ = rootCmd.MarkPersistentFlagRequired("data-dir")

	rootCmd.PersistentFlags().StringVar(&flagPebbleDir, "pebble-dir", "/var/flow/data/protocol-pebble", "directory to the pebble dababase")

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.AutomaticEnv()
}
