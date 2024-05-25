package sikafka_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wonksing/si/v2/sikafka"
)

func TestProducer_Produce(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	producer, err := sikafka.DefaultSyncProducer([]string{"testkafkahost:9092"})
	require.Nil(t, err)
	defer producer.Close()

	sp := sikafka.NewSyncProducer(producer, "tp-test-15")
	p, o, err := sp.Produce([]byte("10123"), []byte("asdf"))
	require.Nil(t, err)
	fmt.Println(p, o)
}

func TestProducer_ProduceWithTopic(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}
	producer, err := sikafka.DefaultSyncProducer([]string{"testkafkahost:9092"})
	require.Nil(t, err)
	defer producer.Close()

	sp := sikafka.NewSyncProducer(producer, "tp-test-15")
	p, o, err := sp.ProduceWithTopic("tp-test", []byte("10123"), []byte("asdf"))
	require.Nil(t, err)
	fmt.Println(p, o)
}
