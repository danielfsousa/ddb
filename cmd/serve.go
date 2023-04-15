package cmd

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/danielfsousa/ddb/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cli = serveCli{}

var serveCmd = &cobra.Command{
	Use:     "serve",
	Short:   "Starts the ddb server",
	PreRunE: cli.setupConfig,
	RunE:    cli.run,
}

func init() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dataDir := path.Join(homedir, ".ddb", "data")

	def := server.NewDefaultConfig()
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringP("host", "H", def.Host, "Host address to bind GRPC server on.")
	serveCmd.Flags().IntP("port", "p", def.Port, "Port for RPC clients connections.")
	serveCmd.Flags().StringP("data-dir", "d", dataDir, "Directory to store database internal data.")

	err = viper.BindPFlags(serveCmd.Flags())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type serveConfig struct {
	Server *server.Config
}

type serveCli struct {
	config *serveConfig
}

func (s *serveCli) setupConfig(_ *cobra.Command, _ []string) error {
	s.config = &serveConfig{Server: &server.Config{}}
	s.config.Server.Host = viper.GetString("host")
	s.config.Server.Port = viper.GetInt("port")
	return nil
}

func (s *serveCli) run(_ *cobra.Command, _ []string) error {
	srv := server.New(s.config.Server)
	if err := srv.Start(); err != nil {
		return err
	}
	runtime.Goexit()
	return nil
}
