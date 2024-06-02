package sigrpc

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "github.com/wonksing/si/v2/sigrpc/tests/protos"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestClient(t *testing.T) {
	serverAddr := "dns:///localhost:50051"
	lis, err := TcpListener(":50051")
	require.Nil(t, err)

	server, err := NewServer(lis, WithDefaultKeepAliveEnforcement(), WithDefaultKeepAlive())
	require.Nil(t, err)
	defer server.Close()

	pb.RegisterStudentServer(server, &studentGrpcServer{})

	go func() {
		if err := server.Start(); err != nil {
			log.Println(err.Error())
		}
	}()

	t.Run("succeed", func(t *testing.T) {

		client, err := NewClient(serverAddr,
			WithDefaultKeepAliveParams(),
			WithInsecureTransportCreds())
		if err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}
		defer client.Close()

		c := pb.NewStudentClient(client)
		res, err := c.Read(context.Background(), &pb.StudentRequest{
			Name: "wonka",
		})
		assert.Nil(t, err)

		require.EqualValues(t, 200, res.Status)
	})
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
