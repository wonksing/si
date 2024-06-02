package sigrpc

import (
	"crypto/tls"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
)

type Server struct {
	Svr      *grpc.Server
	listener *net.Listener
}

// NewServer
func NewServer(listenAddr string,
	enforcementPolicyUse bool, enforcementPolicyMinTime int, enforcementPolicyPermitWithoutStream bool,
	certPem, certKey string,
	keepAliveMaxConnIdle int, keepAliveMaxConnAge int, keepAliveMaxConnAgeGrace int, keepAliveTime int, keepAliveTimeout int,
	healthCheckUse bool,
	opt ...grpc.ServerOption) (*Server, error) {

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}

	opts := []grpc.ServerOption{}
	if certPem != "" && certKey != "" {
		cert, err := tls.LoadX509KeyPair(certPem, certKey)
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.Creds(credentials.NewServerTLSFromCert(&cert)))
	}
	if enforcementPolicyUse {
		kaep := keepalive.EnforcementPolicy{
			MinTime:             time.Duration(enforcementPolicyMinTime) * time.Second,
			PermitWithoutStream: enforcementPolicyPermitWithoutStream,
		}
		opts = append(opts, grpc.KeepaliveEnforcementPolicy(kaep))
	}
	kasp := keepalive.ServerParameters{
		MaxConnectionIdle:     time.Duration(keepAliveMaxConnIdle) * time.Second,
		MaxConnectionAge:      time.Duration(keepAliveMaxConnAge) * time.Second,
		MaxConnectionAgeGrace: time.Duration(keepAliveMaxConnAgeGrace) * time.Second,
		Time:                  time.Duration(keepAliveTime) * time.Second,
		Timeout:               time.Duration(keepAliveTimeout) * time.Second,
	}
	opts = append(opts, grpc.KeepaliveParams(kasp))
	// if apmActive {
	// 	opts = append(opts, grpc.UnaryInterceptor(apmgrpc.NewUnaryServerInterceptor()))
	// 	opts = append(opts, grpc.StreamInterceptor(apmgrpc.NewStreamServerInterceptor()))
	// }
	opts = append(opts, opt...)
	svr := grpc.NewServer(opts...)
	if healthCheckUse {
		healthCheck := health.NewServer()
		healthpb.RegisterHealthServer(svr, healthCheck)
	}

	return &Server{
		Svr:      svr,
		listener: &listener,
	}, nil
}

func (s *Server) Serve(lis net.Listener) error {
	return s.Svr.Serve(lis)
}
func (s *Server) Start() error {
	return s.Svr.Serve(*s.listener)
}
func (s *Server) Stop() error {
	s.Svr.GracefulStop()
	return nil
}
func (s *Server) Close() error {
	s.Svr.GracefulStop()
	return nil
}
