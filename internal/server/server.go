package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/bufbuild/connect-go"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	ddbv1 "github.com/danielfsousa/ddb/gen/ddb/v1"
	"github.com/danielfsousa/ddb/gen/ddb/v1/ddbv1connect"
)

var ErrEmptyKey = errors.New("key cannot be empty")

// Server implements the DdbService API.
type Server struct {
	ddbv1connect.UnimplementedDdbServiceHandler
	config     *Config
	httpServer *http2.Server
	mu         sync.Mutex
	data       map[string][]byte
}

var _ ddbv1connect.DdbServiceHandler = (*Server)(nil)

// New will create a new Server.
func New(config *Config) *Server {
	return &Server{
		config: config,
		data:   make(map[string][]byte),
	}
}

// Start will start the Server and block until it is signaled to stop.
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	mux := http.NewServeMux()
	path, handler := ddbv1connect.NewDdbServiceHandler(s)
	mux.Handle(path, handler)
	s.httpServer = &http2.Server{}

	// sigc := make(chan os.Signal, 1)
	// signal.Notify(sigc, os.Interrupt)
	// go func() {
	// 	<-sigc
	// 	// TODO: graceful shutdown
	// }()

	fmt.Println("Server listening on", addr)
	return http.ListenAndServe( //nolint:gosec // TODO: fix this
		addr,
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, s.httpServer),
	)
}

// Stop will attempt to gracefully shut the Server down by signaling the stop.
func (s *Server) Stop() error {
	if s.httpServer != nil {
		fmt.Println("should stop server")
		// 	// TODO: graceful shutdown
		// s.httpServer.Stop()
	}
	return nil
}

// Get will return the value for the given key.
func (s *Server) Get(
	_ context.Context,
	req *connect.Request[ddbv1.GetRequest],
) (*connect.Response[ddbv1.GetResponse], error) {
	key := req.Msg.GetKey()
	if err := validateKey(key); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.data[key]
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("could not find key %q", key))
	}

	return connect.NewResponse(&ddbv1.GetResponse{Key: key, Value: value}), nil
}

// Set will set the value for the given key.
func (s *Server) Set(
	_ context.Context,
	req *connect.Request[ddbv1.SetRequest],
) (*connect.Response[ddbv1.SetResponse], error) {
	key := req.Msg.GetKey()
	value := req.Msg.GetValue()
	if err := validateKey(key); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value

	return connect.NewResponse(&ddbv1.SetResponse{}), nil
}

func validateKey(key string) error {
	if key == "" {
		return ErrEmptyKey
	}
	return nil
}
