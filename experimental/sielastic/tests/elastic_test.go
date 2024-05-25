package sielastic_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wonksing/si/v2/experimental/sielastic"
	"github.com/wonksing/si/v2/tests/testmodels"
)

func TestElasticClient_IndexDocument(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	require.NotNil(t, client)

	ec := sielastic.NewClient(client)
	// require.Nil(t, err)

	d := testmodels.Document{
		Name:      "my name is wonk",
		ID:        1040289,
		Timestamp: time.Now(),
	}
	body := d.String()
	res, err := ec.IndexDocument(context.Background(), "idx-test", []byte(body))
	require.Nil(t, err)

	fmt.Println(res)
}

func TestElasticClient_SearchDocuments(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	require.NotNil(t, client)

	ec := sielastic.NewClient(client)
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"name": "name",
			},
		},
	}

	res := make(map[string]interface{})
	err := ec.SearchDocuments(context.Background(), "idx-test", query, &res)
	require.Nil(t, err)
	fmt.Println(res)
}

func TestElasticClient_SearchDocuments_Fail(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	require.NotNil(t, client)

	ec := sielastic.NewClient(client)
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"name": "name",
			},
		},
		"query2": map[string]interface{}{
			"match": map[string]interface{}{
				"name": "name",
			},
		},
	}

	res := make(map[string]interface{})
	err := ec.SearchDocuments(context.Background(), "idx-test", query, &res)
	require.NotNil(t, err)
	fmt.Println(res)
}
