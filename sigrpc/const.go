package sigrpc

import "time"

const (
	// both client and server
	defaultPermitWithoutStream = true

	// client
	defaultKeepAliveTime    = 10 * time.Second
	defaultKeepAliveTimeout = 3 * time.Second
	defaultDialTimeout      = 6 * time.Second
	defaultDialBlock        = false

	defaultServiceConfig = `{
		"loadBalancingConfig": [{"round_robin":{}}],
		"methodConfig": [{
			"name": [],
			"waitForReady": true,
			"retryPolicy": {
				"maxAttempts": 3,
				"initialBackoff": ".1s",
				"maxBackoff": "3s",
				"backoffMultiplier": 1.5,
				"retryableStatusCodes": [ "UNAVAILABLE", "RESOURCE_EXHASTED" ]
			}
		}]
	}`

	// server
	defaultMinTime               = 5 * time.Second
	defaultMaxConnectionIdle     = defaultKeepAliveTime + defaultKeepAliveTimeout
	defaultMaxConnectionAge      = defaultMaxConnectionIdle + (10 * time.Second)
	defaultMaxConnectionAgeGrace = 5 * time.Second
	defaultTime                  = 6 * time.Second
	defaultTimeout               = 3 * time.Second
)
