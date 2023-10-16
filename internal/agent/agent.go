// package agent sets up the DDB server, logger and membership components.
package agent

import (
	"fmt"
	"net"
	"sync"

	"github.com/danielfsousa/ddb"
	"github.com/danielfsousa/ddb/internal/discovery"
	"github.com/danielfsousa/ddb/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Agent struct {
	Config *Config

	database   *ddb.Ddb
	server     *server.Server
	membership *discovery.Membership

	shutdown     bool
	shutdowns    chan struct{}
	shutdownLock sync.Mutex
	logger       *zerolog.Logger
}

func (c *Config) RPCAddr() (string, error) {
	host, _, err := net.SplitHostPort(c.BindAddr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", host, c.RPCPort), nil
}

func New(config *Config) (*Agent, error) {
	logger := log.With().Str("component", "agent").Logger()
	agent := &Agent{
		Config:    config,
		shutdowns: make(chan struct{}),
		logger:    &logger,
	}

	if err := agent.setupDatabase(); err != nil {
		return nil, err
	}
	if err := agent.setupServer(); err != nil {
		return nil, err
	}
	if err := agent.setupMembership(); err != nil {
		return nil, err
	}

	return agent, nil
}

func (a *Agent) setupDatabase() (err error) {
	a.database, err = ddb.Open(a.Config.DataDir)
	return err
}

func (a *Agent) setupServer() error {
	host, _, err := net.SplitHostPort(a.Config.BindAddr)
	if err != nil {
		return err
	}

	a.server = server.New(&server.Config{
		Host: host,
		Port: a.Config.RPCPort,
		Ddb:  a.database,
	})
	go func() {
		if err := a.server.Start(); err != nil {
			a.logger.Fatal().Err(err).Msg("failed to start server")
			_ = a.Shutdown()
		}
	}()
	return nil
}

type tempDiscoveryHandler struct {
	logger *zerolog.Logger
}

func newTempDiscoveryHandler() *tempDiscoveryHandler {
	logger := log.With().Str("component", "temp-discovery-handler").Logger()
	return &tempDiscoveryHandler{logger: &logger}
}

func (t tempDiscoveryHandler) Join(name, addr string) error {
	t.logger.Debug().Str("name", name).Str("addr", addr).Msg("join event")
	return nil
}

func (t tempDiscoveryHandler) Leave(name string) error {
	t.logger.Debug().Str("name", name).Msg("leave event")
	return nil
}

func (a *Agent) setupMembership() error {
	rpcAddr, err := a.Config.RPCAddr()
	if err != nil {
		return err
	}
	tempHandler := newTempDiscoveryHandler()
	a.membership, err = discovery.New(tempHandler, discovery.Config{
		NodeName: a.Config.NodeName,
		BindAddr: a.Config.BindAddr,
		Tags: map[string]string{
			"rpc_addr": rpcAddr,
		},
		StartJoinAddrs: a.Config.StartJoinAddrs,
	})
	return err
}

func (a *Agent) Shutdown() error {
	a.shutdownLock.Lock()
	defer a.shutdownLock.Unlock()
	if a.shutdown {
		return nil
	}
	a.shutdown = true
	a.logger.Info().Msg("shutting down")
	close(a.shutdowns)

	shutdown := []func() error{
		a.membership.Leave,
		a.server.Stop,
		a.database.Close,
	}
	for _, fn := range shutdown {
		if err := fn(); err != nil {
			a.logger.Error().Err(err).Msg("failed to shutdown gracefully")
			return err
		}
	}
	return nil
}
