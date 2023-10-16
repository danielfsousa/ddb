package main

import (
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/danielfsousa/ddb/internal/agent"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	logger  *zerolog.Logger
)

func main() {
	logger = setupGlobalLogger()
	cli := ddbServerCli{}
	cmd := &cobra.Command{
		Use:     "ddb-server",
		Short:   "A distributed key-value database",
		Version: "0.1.0",
		PreRun:  cli.setup,
		RunE:    cli.run,
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to get user home directory")
	}

	def := agent.NewDefaultConfig()

	cobra.OnInitialize(setupConfigFile)
	cmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.ddb/config.yaml)")
	// TODO: split bind-addr into host and serf-port?
	cmd.Flags().StringP("bind-addr", "a", def.BindAddr, "Address to bind Serf on.")
	cmd.Flags().StringP("data-dir", "d", path.Join(homeDir, ".ddb", "data"), "Directory to store database internal data.")
	cmd.Flags().IntP("rpc-port", "p", def.RPCPort, "Port for RPC clients (and Raft) connections.")

	err = viper.BindPFlags(cmd.Flags())
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to parse flags")
	}

	if err := cmd.Execute(); err != nil {
		logger.Fatal().Err(err).Msg("failed to execute command")
	}
}

func setupGlobalLogger() *zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	logger := log.With().Str("component", "main").Logger()
	return &logger
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
		logger.Info().Msgf("using config file: %s", viper.ConfigFileUsed())
	}
}

type ddbServerCli struct {
	config *agent.Config
	logger *zerolog.Logger
}

func (cli *ddbServerCli) setup(cmd *cobra.Command, args []string) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	logger := log.With().Str("component", "main").Logger()
	cli.logger = &logger
	cli.config = &agent.Config{
		DataDir:        viper.GetString("data-dir"),
		NodeName:       viper.GetString("node-name"),
		BindAddr:       viper.GetString("bind-addr"),
		RPCPort:        viper.GetInt("rpc-port"),
		StartJoinAddrs: viper.GetStringSlice("start-join-addrs"),
		Bootstrap:      viper.GetBool("bootstrap"),
	}
}

func (cli *ddbServerCli) run(cmd *cobra.Command, args []string) error {
	a, err := agent.New(cli.config)
	if err != nil {
		return err
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc
	return a.Shutdown()
}
