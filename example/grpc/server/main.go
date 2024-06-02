package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/wonksing/si/v2/example/grpc/protos"
	"github.com/wonksing/si/v2/sigrpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {

	var err error
	serverAddr := ":50051"

	listener, _ := sigrpc.TcpListener(serverAddr)
	server, err := sigrpc.NewServer(listener,
		sigrpc.WithDefaultKeepAliveEnforcement(),
		sigrpc.WithDefaultKeepAlive())

	// build server
	// enforcementPolicyUse := true
	// enforcementPolicyMinTime := 15
	// enforcementPolicyPermitWithoutStream := true
	// certPem := ""
	// certKey := ""
	// keepAliveMaxConnIdle := 300
	// keepAliveMaxConnAge := 300
	// keepAliveMaxConnAgeGrace := 6
	// keepAliveTime := 60
	// keepAliveTimeout := 1
	// healthCheckUse := true
	// listener, _ := sigrpc.TcpListener(serverAddr)
	// server, err := sigrpc.NewServer(listener,
	// 	enforcementPolicyUse, enforcementPolicyMinTime, enforcementPolicyPermitWithoutStream,
	// 	certPem, certKey,
	// 	keepAliveMaxConnIdle, keepAliveMaxConnAge, keepAliveMaxConnAgeGrace, keepAliveTime, keepAliveTimeout,
	// 	healthCheckUse)
	if err != nil {
		os.Exit(1)
	}
	pb.RegisterStudentServer(server, &studentGrpcServer{})

	if err := server.Start(); err != nil {
		log.Println(err.Error())
	}
}

type studentGrpcServer struct {
	pb.StudentServer
}

func (d *studentGrpcServer) Read(ctx context.Context, in *pb.StudentRequest) (*pb.StudentReply, error) {
	docs := make([]*pb.StudentEntity, 0)
	docs = append(docs, &pb.StudentEntity{
		Name:        "wonk",
		Age:         10,
		DateTime:    timestamppb.New(time.Now()),
		DoubleValue: 10.1,
	})

	var count int64 = 1
	rep := pb.StudentReply{
		Status:    200,
		Documents: docs,
		Count:     &count,
	}

	return &rep, nil
}
