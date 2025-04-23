package domain

import (
	"context"
)

type StreamResponse interface {
	Content() (string, bool, error)
}

type Agent interface {
	GetResponse(ctx context.Context, messages []Message) (string, error)
	GetStreamResponse(ctx context.Context, messages []Message) (StreamResponse, error)
}
