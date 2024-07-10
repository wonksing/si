package sigrpc

import (
	"context"
	"crypto/tls"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func WithInsecureTransportCreds() grpc.DialOption {
	return grpc.WithTransportCredentials(insecure.NewCredentials())
}

func WithTLSConfigTransportCreds(c *tls.Config) grpc.DialOption {
	// &tls.Config{} for example
	return grpc.WithTransportCredentials(credentials.NewTLS(c))
}
func WithDefaultKeepAliveParams() grpc.DialOption {
	kacp := keepalive.ClientParameters{
		Time:                defaultKeepAliveTime,
		Timeout:             defaultKeepAliveTimeout,
		PermitWithoutStream: defaultPermitWithoutStream,
	}
	return grpc.WithKeepaliveParams(kacp)
}
func WithKeepAliveParams(keepAliveTime, keepAliveTimeout time.Duration, permitWithoutStream bool) grpc.DialOption {
	kacp := keepalive.ClientParameters{
		Time:                keepAliveTime,
		Timeout:             keepAliveTimeout,
		PermitWithoutStream: permitWithoutStream,
	}
	return grpc.WithKeepaliveParams(kacp)
}

func TransportCredentialsOption(certPemFile string, serverNameOverride string) (grpc.DialOption, error) {
	creds, err := credentials.NewClientTLSFromFile(certPemFile, serverNameOverride)
	if err != nil {
		return nil, err
	}
	return grpc.WithTransportCredentials(creds), nil
}

func WithDefaultServiceConfig(defaultServiceConfig string) grpc.DialOption {
	return grpc.WithDefaultServiceConfig(defaultServiceConfig)
}

func WithDefaultDialBlock() grpc.DialOption {
	return WithDialBlock(defaultDialBlock)
}
func WithDialBlock(dialBlock bool) grpc.DialOption {
	if dialBlock {
		return grpc.WithBlock()
	}
	return &grpc.EmptyDialOption{}
}

func NewClient(address string, opts ...grpc.DialOption) (*Client, error) {
	return NewClientWithDialTimeout(address, defaultDialTimeout, opts...)
}

func NewClientWithDialTimeout(address string, dialTimeout time.Duration, opts ...grpc.DialOption) (*Client, error) {
	defaultOpts := []grpc.DialOption{}

	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()

	defaultOpts = append(defaultOpts, opts...)
	conn, err := grpc.DialContext(ctx, address,
		defaultOpts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		ClientConn: conn,
	}, nil
}

type Client struct {
	*grpc.ClientConn
}

// // NewClient returns Client
// //
// // Deprecated
// func NewClient(
// 	addrs, resolverScheme, resolverSchemeName string,
// 	keepAliveTime, keepAliveTimeout int, keepAlivePermitWithoutStream bool,
// 	caCertPem, certServername string,
// 	defaultServiceConfig string, dialBlock bool, dialTimeoutSecond int,
// ) (*Client, error) {

// 	resolver.Register(&grpcResolverBuilder{
// 		scheme:      resolverScheme,
// 		serviceName: resolverSchemeName,
// 		addrs:       strings.Split(addrs, ","),
// 	})
// 	kacp := keepalive.ClientParameters{
// 		Time:                time.Duration(keepAliveTime) * time.Second,
// 		Timeout:             time.Duration(keepAliveTimeout) * time.Second,
// 		PermitWithoutStream: keepAlivePermitWithoutStream,
// 	}
// 	opts := []grpc.DialOption{
// 		grpc.WithKeepaliveParams(kacp),
// 	}
// 	if caCertPem != "" && certServername != "" {
// 		creds, err := credentials.NewClientTLSFromFile(caCertPem, certServername)
// 		if err != nil {
// 			return nil, err
// 		}
// 		opts = append(opts, grpc.WithTransportCredentials(creds))
// 	}
// 	if defaultServiceConfig != "" {
// 		opts = append(opts, grpc.WithDefaultServiceConfig(defaultServiceConfig))
// 	}
// 	if dialBlock {
// 		opts = append(opts, grpc.WithBlock())
// 	}
// 	var dialTimeout time.Duration
// 	if dialTimeoutSecond == 0 {
// 		dialTimeout = 12 * time.Second
// 	} else {
// 		dialTimeout = time.Duration(dialTimeoutSecond) * time.Second
// 	}
// 	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
// 	defer cancel()
// 	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:///%s", resolverScheme, resolverSchemeName),
// 		opts...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Client{
// 		ClientConn: conn,
// 	}, nil
// }
