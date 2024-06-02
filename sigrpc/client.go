package sigrpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/resolver"
)

type Client struct {
	*grpc.ClientConn
}

// NewClient
func NewClient(
	addrs, resolverScheme, resolverSchemeName string,
	keepAliveTime, keepAliveTimeout int, keepAlivePermitWithoutStream bool,
	caCertPem, certServername string,
	defaultServiceConfig string, dialBlock bool, dialTimeoutSecond int,
) (*Client, error) {

	resolver.Register(&grpcResolverBuilder{
		scheme:      resolverScheme,
		serviceName: resolverSchemeName,
		addrs:       strings.Split(addrs, ","),
	})
	kacp := keepalive.ClientParameters{
		Time:                time.Duration(keepAliveTime) * time.Second,
		Timeout:             time.Duration(keepAliveTimeout) * time.Second,
		PermitWithoutStream: keepAlivePermitWithoutStream,
	}
	opts := []grpc.DialOption{
		grpc.WithKeepaliveParams(kacp),
	}
	if caCertPem != "" && certServername != "" {
		creds, err := credentials.NewClientTLSFromFile(caCertPem, certServername)
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}
	if defaultServiceConfig != "" {
		opts = append(opts, grpc.WithDefaultServiceConfig(defaultServiceConfig))
	}
	if dialBlock {
		opts = append(opts, grpc.WithBlock())
	}
	var dialTimeout time.Duration
	if dialTimeoutSecond == 0 {
		dialTimeout = 12 * time.Second
	} else {
		dialTimeout = time.Duration(dialTimeoutSecond) * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:///%s", resolverScheme, resolverSchemeName),
		opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		ClientConn: conn,
	}, nil
}
