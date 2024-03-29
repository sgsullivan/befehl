package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Config *viper.Viper

var RootCmd = &cobra.Command{
	Use: "befehl",
	Long: `ausführen willkürliche Befehle über ssh in Masse

Dieses Werkzeug sollte mit Vorsicht verwendet werden; gegeben das Macht angeboten.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		color.Red("A fatal error has occurred:")
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	thisViper := viper.New()
	thisViper.SetConfigName(".befehl")
	configDir := os.Getenv("HOME")
	if os.Getenv("BEFEHL_CONFIG_DIR") != "" {
		configDir = os.Getenv("BEFEHL_CONFIG_DIR")
	}
	thisViper.AddConfigPath(configDir)
	thisViper.SetEnvPrefix("befehl")
	thisViper.AutomaticEnv()
	thisViper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := thisViper.ReadInConfig(); err != nil {
		fmt.Printf("Failed reading config (using defaults) [%s]: [%s]\n", thisViper.ConfigFileUsed(), err)
	}

	Config = thisViper
}
