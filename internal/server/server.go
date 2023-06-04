package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bufbuild/connect-go"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/danielfsousa/ddb"
	ddbv1 "github.com/danielfsousa/ddb/gen/ddb/v1"
	"github.com/danielfsousa/ddb/gen/ddb/v1/ddbv1connect"
)

// Server implements the DdbService API.
type Server struct {
	ddbv1connect.UnimplementedDdbServiceHandler
	httpServer *http2.Server
	*Config
}

var _ ddbv1connect.DdbServiceHandler = (*Server)(nil)

// New will create a new Server.
func New(config *Config) *Server {
	return &Server{Config: config}
}

// Start will start the Server and block until it is signaled to stop.
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
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

// Has will return true if the given key exists in the database.
func (s *Server) Has(
	_ context.Context,
	req *connect.Request[ddbv1.HasRequest],
) (*connect.Response[ddbv1.HasResponse], error) {
	key := req.Msg.GetKey()
	if err := validateKey(key); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	exists := s.Ddb.Has(key)

	return connect.NewResponse(&ddbv1.HasResponse{Key: key, Exists: exists}), nil
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

	value, err := s.Ddb.Get(key)
	if err != nil {
		if err == ddb.ErrKeyNotFound {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
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

	err := s.Ddb.Set(key, value)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&ddbv1.SetResponse{}), nil
}

func (s *Server) Delete(
	_ context.Context,
	req *connect.Request[ddbv1.DeleteRequest],
) (*connect.Response[ddbv1.DeleteResponse], error) {
	key := req.Msg.GetKey()
	if err := validateKey(key); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if !s.Ddb.Has(key) {
		return nil, connect.NewError(connect.CodeNotFound, ddb.ErrKeyNotFound)
	}

	err := s.Ddb.Delete(key)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	return connect.NewResponse(&ddbv1.DeleteResponse{}), nil
}

func validateKey(key string) error {
	if key == "" {
		return ddb.ErrKeyEmpty
	}
	return nil
}
