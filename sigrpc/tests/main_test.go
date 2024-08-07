package sigrpc_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/wonksing/si/v2/sigrpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/wonksing/si/v2/sigrpc/tests/protos"
)

var (
	onlinetest, _ = strconv.ParseBool(os.Getenv("ONLINE_TEST"))

	server     *sigrpc.Server
	serverAddr = ":60000"

	client *sigrpc.Client
)

func setup() error {
	var err error

	// build server
	// enforcementPolicyUse := true
	// enforcementPolicyMinTime := 15
	// enforcementPolicyPermitWithoutStream := true
	// certPem := "./certs/server.crt"
	// certKey := "./certs/server.key"
	// keepAliveMaxConnIdle := 300
	// keepAliveMaxConnAge := 300
	// keepAliveMaxConnAgeGrace := 6
	// keepAliveTime := 60
	// keepAliveTimeout := 1
	// healthCheckUse := true

	// listener, _ := sigrpc.TcpListener(serverAddr)
	// server, err = sigrpc.NewServer(listener,
	// 	enforcementPolicyUse, enforcementPolicyMinTime, enforcementPolicyPermitWithoutStream,
	// 	certPem, certKey,
	// 	keepAliveMaxConnIdle, keepAliveMaxConnAge, keepAliveMaxConnAgeGrace, keepAliveTime, keepAliveTimeout,
	// 	healthCheckUse)
	// if err != nil {
	// 	return err
	// }

	lis, _ := sigrpc.TcpListener(serverAddr)
	server, err = sigrpc.NewServer(lis)
	if err != nil {
		return err
	}
	pb.RegisterStudentServer(server, &studentGrpcServer{})

	go func() {
		server.Start()
	}()

	// build client
	client, err = sigrpc.NewClient(serverAddr,
		sigrpc.WithInsecureTransportCreds())
	if err != nil {
		return err
	}

	if onlinetest {
	}

	return nil
}

func shutdown() {
	if server != nil {
		server.Close()
	}
	if client != nil {
		client.Close()
	}
}

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		fmt.Println(err)
		shutdown()
		os.Exit(1)
	}

	exitCode := m.Run()

	shutdown()
	os.Exit(exitCode)
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
