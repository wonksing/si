package sikafka

import (
	"github.com/IBM/sarama"
	"github.com/eapache/go-resiliency/breaker"
)

func isRetryableError(err error) bool {
	switch err {
	case sarama.ErrBrokerNotAvailable,
		sarama.ErrLeaderNotAvailable,
		sarama.ErrReplicaNotAvailable,
		sarama.ErrRequestTimedOut,
		sarama.ErrNotEnoughReplicas,
		// sarama.ErrNotEnoughReplicasAfterAppend, // "kafka server: Messages are written to the log, but to fewer in-sync replicas than required"
		// sarama.ErrNetworkException, // "kafka server: The server disconnected before a response was received"
		sarama.ErrOutOfBrokers,
		sarama.ErrOutOfOrderSequenceNumber,
		sarama.ErrNotController,
		sarama.ErrNotLeaderForPartition,
		breaker.ErrBreakerOpen:
		return true
	default:
		return false
	}
}
