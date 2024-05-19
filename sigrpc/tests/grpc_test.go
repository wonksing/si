package sigrpc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	pb "github.com/wonksing/si/v2/sigrpc/tests/protos"
)

func TestXxx(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	c := pb.NewStudentClient(client)
	rep, err := c.Read(context.Background(), &pb.StudentRequest{
		Name: "wonka",
	})
	assert.Nil(t, err)

	fmt.Println(rep.String())
}
