// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: ddb/v1/ddb.proto

package ddbv1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1 "github.com/danielfsousa/ddb/gen/ddb/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// DdbServiceName is the fully-qualified name of the DdbService service.
	DdbServiceName = "ddb.v1.DdbService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// DdbServiceHasProcedure is the fully-qualified name of the DdbService's Has RPC.
	DdbServiceHasProcedure = "/ddb.v1.DdbService/Has"
	// DdbServiceGetProcedure is the fully-qualified name of the DdbService's Get RPC.
	DdbServiceGetProcedure = "/ddb.v1.DdbService/Get"
	// DdbServiceSetProcedure is the fully-qualified name of the DdbService's Set RPC.
	DdbServiceSetProcedure = "/ddb.v1.DdbService/Set"
	// DdbServiceDeleteProcedure is the fully-qualified name of the DdbService's Delete RPC.
	DdbServiceDeleteProcedure = "/ddb.v1.DdbService/Delete"
)

// DdbServiceClient is a client for the ddb.v1.DdbService service.
type DdbServiceClient interface {
	Has(context.Context, *connect_go.Request[v1.HasRequest]) (*connect_go.Response[v1.HasResponse], error)
	Get(context.Context, *connect_go.Request[v1.GetRequest]) (*connect_go.Response[v1.GetResponse], error)
	Set(context.Context, *connect_go.Request[v1.SetRequest]) (*connect_go.Response[v1.SetResponse], error)
	Delete(context.Context, *connect_go.Request[v1.DeleteRequest]) (*connect_go.Response[v1.DeleteResponse], error)
}

// NewDdbServiceClient constructs a client for the ddb.v1.DdbService service. By default, it uses
// the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewDdbServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) DdbServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &ddbServiceClient{
		has: connect_go.NewClient[v1.HasRequest, v1.HasResponse](
			httpClient,
			baseURL+DdbServiceHasProcedure,
			opts...,
		),
		get: connect_go.NewClient[v1.GetRequest, v1.GetResponse](
			httpClient,
			baseURL+DdbServiceGetProcedure,
			opts...,
		),
		set: connect_go.NewClient[v1.SetRequest, v1.SetResponse](
			httpClient,
			baseURL+DdbServiceSetProcedure,
			opts...,
		),
		delete: connect_go.NewClient[v1.DeleteRequest, v1.DeleteResponse](
			httpClient,
			baseURL+DdbServiceDeleteProcedure,
			opts...,
		),
	}
}

// ddbServiceClient implements DdbServiceClient.
type ddbServiceClient struct {
	has    *connect_go.Client[v1.HasRequest, v1.HasResponse]
	get    *connect_go.Client[v1.GetRequest, v1.GetResponse]
	set    *connect_go.Client[v1.SetRequest, v1.SetResponse]
	delete *connect_go.Client[v1.DeleteRequest, v1.DeleteResponse]
}

// Has calls ddb.v1.DdbService.Has.
func (c *ddbServiceClient) Has(ctx context.Context, req *connect_go.Request[v1.HasRequest]) (*connect_go.Response[v1.HasResponse], error) {
	return c.has.CallUnary(ctx, req)
}

// Get calls ddb.v1.DdbService.Get.
func (c *ddbServiceClient) Get(ctx context.Context, req *connect_go.Request[v1.GetRequest]) (*connect_go.Response[v1.GetResponse], error) {
	return c.get.CallUnary(ctx, req)
}

// Set calls ddb.v1.DdbService.Set.
func (c *ddbServiceClient) Set(ctx context.Context, req *connect_go.Request[v1.SetRequest]) (*connect_go.Response[v1.SetResponse], error) {
	return c.set.CallUnary(ctx, req)
}

// Delete calls ddb.v1.DdbService.Delete.
func (c *ddbServiceClient) Delete(ctx context.Context, req *connect_go.Request[v1.DeleteRequest]) (*connect_go.Response[v1.DeleteResponse], error) {
	return c.delete.CallUnary(ctx, req)
}

// DdbServiceHandler is an implementation of the ddb.v1.DdbService service.
type DdbServiceHandler interface {
	Has(context.Context, *connect_go.Request[v1.HasRequest]) (*connect_go.Response[v1.HasResponse], error)
	Get(context.Context, *connect_go.Request[v1.GetRequest]) (*connect_go.Response[v1.GetResponse], error)
	Set(context.Context, *connect_go.Request[v1.SetRequest]) (*connect_go.Response[v1.SetResponse], error)
	Delete(context.Context, *connect_go.Request[v1.DeleteRequest]) (*connect_go.Response[v1.DeleteResponse], error)
}

// NewDdbServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewDdbServiceHandler(svc DdbServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle(DdbServiceHasProcedure, connect_go.NewUnaryHandler(
		DdbServiceHasProcedure,
		svc.Has,
		opts...,
	))
	mux.Handle(DdbServiceGetProcedure, connect_go.NewUnaryHandler(
		DdbServiceGetProcedure,
		svc.Get,
		opts...,
	))
	mux.Handle(DdbServiceSetProcedure, connect_go.NewUnaryHandler(
		DdbServiceSetProcedure,
		svc.Set,
		opts...,
	))
	mux.Handle(DdbServiceDeleteProcedure, connect_go.NewUnaryHandler(
		DdbServiceDeleteProcedure,
		svc.Delete,
		opts...,
	))
	return "/ddb.v1.DdbService/", mux
}

// UnimplementedDdbServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedDdbServiceHandler struct{}

func (UnimplementedDdbServiceHandler) Has(context.Context, *connect_go.Request[v1.HasRequest]) (*connect_go.Response[v1.HasResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("ddb.v1.DdbService.Has is not implemented"))
}

func (UnimplementedDdbServiceHandler) Get(context.Context, *connect_go.Request[v1.GetRequest]) (*connect_go.Response[v1.GetResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("ddb.v1.DdbService.Get is not implemented"))
}

func (UnimplementedDdbServiceHandler) Set(context.Context, *connect_go.Request[v1.SetRequest]) (*connect_go.Response[v1.SetResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("ddb.v1.DdbService.Set is not implemented"))
}

func (UnimplementedDdbServiceHandler) Delete(context.Context, *connect_go.Request[v1.DeleteRequest]) (*connect_go.Response[v1.DeleteResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("ddb.v1.DdbService.Delete is not implemented"))
}
