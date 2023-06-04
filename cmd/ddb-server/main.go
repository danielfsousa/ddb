package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/danielfsousa/ddb"
	"github.com/danielfsousa/ddb/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func main() {
	cmd := &cobra.Command{
		Use:     "ddb-server",
		Short:   "A distributed key-value database",
		Version: "0.1.0",
		RunE:    run,
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	def := server.NewDefaultConfig()

	cobra.OnInitialize(setupConfigFile)
	cmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.ddb/config.yaml)")
	cmd.Flags().StringP("host", "H", def.Host, "Host address to bind GRPC server on.")
	cmd.Flags().IntP("port", "p", def.Port, "Port for RPC clients connections.")
	cmd.Flags().StringP("data-dir", "d", path.Join(homeDir, ".ddb", "data"), "Directory to store database internal data.")

	err = viper.BindPFlags(cmd.Flags())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setupConfigFile() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		ddbdir := path.Join(home, ".ddb")
		viper.AddConfigPath(ddbdir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func run(_ *cobra.Command, _ []string) error {
	database, err := ddb.Open(viper.GetString("data-dir"))
	if err != nil {
		return err
	}

	config := &server.Config{}
	config.Ddb = database
	config.Host = viper.GetString("host")
	config.Port = viper.GetInt("port")

	srv := server.New(config)
	if err := srv.Start(); err != nil {
		return err
	}
	runtime.Goexit()
	return nil
}
