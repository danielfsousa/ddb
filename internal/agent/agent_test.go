package agent_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/bufbuild/connect-go"
	ddbv1 "github.com/danielfsousa/ddb/gen/ddb/v1"
	"github.com/danielfsousa/ddb/gen/ddb/v1/ddbv1connect"
	agent "github.com/danielfsousa/ddb/internal/agent"
	"github.com/stretchr/testify/require"
	"github.com/travisjeffery/go-dynaport"
)

func TestAgent(t *testing.T) {
	var agents []*agent.Agent
	for i := 0; i < 3; i++ {
		ports := dynaport.Get(2)
		bindAddr := fmt.Sprintf("%s:%d", "127.0.0.1", ports[0])
		rpcPort := ports[1]

		dataDir := t.TempDir()
		var startJoinAddrs []string
		if i != 0 {
			startJoinAddrs = append(startJoinAddrs, agents[0].Config.BindAddr)
		}

		a, err := agent.New(&agent.Config{
			NodeName:       fmt.Sprintf("node-%d", i),
			StartJoinAddrs: startJoinAddrs,
			BindAddr:       bindAddr,
			RPCPort:        rpcPort,
			DataDir:        dataDir,
			Bootstrap:      i == 0,
		})
		require.NoError(t, err)
		agents = append(agents, a)
	}

	defer func() {
		for _, agent := range agents {
			err := agent.Shutdown()
			require.NoError(t, err)
		}
	}()

	time.Sleep(3 * time.Second)

	leaderClient := client(t, agents[0])
	_, err := leaderClient.Set(
		context.Background(),
		connect.NewRequest(&ddbv1.SetRequest{Key: "foo", Value: []byte("bar")}),
	)
	require.NoError(t, err)

	// TODO: test replication
}

func client(t *testing.T, a *agent.Agent) ddbv1connect.DdbServiceClient {
	addr, err := a.Config.RPCAddr()
	require.NoError(t, err)
	return ddbv1connect.NewDdbServiceClient(
		http.DefaultClient,
		"http://"+addr,
	)
}
