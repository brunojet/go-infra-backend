package adapters

import (
	"context"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

// EndpointFromEnv returns OTLP endpoint configured via OTLP_ENDPOINT or default.
func EndpointFromEnv() string {
	endpoint := os.Getenv("OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4317"
	}
	return endpoint
}

// DialOTLP creates a gRPC ClientConn to the OTLP endpoint using the non-deprecated
// Client API. It waits up to 5s for the connection to become Ready.
func DialOTLP(ctx context.Context, endpoint string) (*grpc.ClientConn, error) {
	// Use NewClient which returns a *ClientConn. It accepts DialOptions.
	cc, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// Kick connection out of idle and wait a short time for READY.
	if cc.GetState() == connectivity.Idle {
		cc.Connect()
	}
	waitCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	for {
		s := cc.GetState()
		if s == connectivity.Ready {
			return cc, nil
		}
		if !cc.WaitForStateChange(waitCtx, s) {
			// timeout or cancelled; still return the client though.
			return cc, nil
		}
	}
}
