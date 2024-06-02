package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/wonksing/si/v2/sigrpc"
	pb "github.com/wonksing/si/v2/sigrpc/tests/protos"
)

func main() {
	client1()
	client2()
}

func client1() {
	// sigrpc.DefaultClient("dns://")

	// build client
	serverAddr := "dns:///localhost:50051"

	client, err := sigrpc.NewClient(serverAddr,
		sigrpc.WithInsecureTransportCreds(),
		sigrpc.WithDefaultKeepAliveParams())
	// client, err := sigrpc.NewClient(serverAddr, resolveScheme, resolveServiceName, keepAliveTime, keepAliveTimeout, keepAlivePermitWithoutStream,
	// 	certPem, certServername, defaultServiceConfig, dialBlock, dialTimeoutSecond)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	defer client.Close()

	c := pb.NewStudentClient(client)
	rep, err := c.Read(context.Background(), &pb.StudentRequest{
		Name: "wonka",
	})
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(rep.String())
}
func client2() {
	// sigrpc.DefaultClient("dns://")

	// build client
	serverAddr := "dns:///fw-grpc.wonksing.com"

	client, err := sigrpc.NewClient(serverAddr,
		sigrpc.WithTLSConfigTransportCreds(&tls.Config{}),
		sigrpc.WithDefaultKeepAliveParams())
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	defer client.Close()

	c := pb.NewStudentClient(client)
	rep, err := c.Read(context.Background(), &pb.StudentRequest{
		Name: "wonka",
	})
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(rep.String())
}
