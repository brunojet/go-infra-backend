package message

// Package message exposes message-related contracts for convenience.

import (
	"context"

	imsg "github.com/brunojet/go-infra-backend/internal/infra/message/adapters"
	"github.com/brunojet/go-infra-backend/pkg/infra/message/contracts"
)

type (
	MessageQueueAdapter = contracts.MessageQueueAdapter
)

// NewLocalS3EventQueue constructs the local S3 adapter. Caller must provide
// an absolute storagePath. Returns the adapter as the exported contract.
func NewLocalS3EventQueue(ctx context.Context, storagePath, filePath string) (MessageQueueAdapter, error) {
	a, err := imsg.NewLocalS3EventQueue(ctx, storagePath, filePath)
	if err != nil {
		return nil, err
	}
	return a, nil
}
