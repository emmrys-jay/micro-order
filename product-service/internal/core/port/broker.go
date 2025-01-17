package port

import (
	"context"

	"go.uber.org/zap"
)

// MessageQueueRepository is an interface for interacting with message queue logic
type MessageQueueRepository interface {
	// Publish pushes messages to the queue
	Publish(ctx context.Context, queue string, msg []byte, headers map[string]any) error
	// Consume retrieves messages from the queue
	Consume(ctx context.Context, queue string, handler func(*zap.Logger, []byte) error)
}
