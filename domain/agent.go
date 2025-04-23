package domain

import (
	"context"
)

type Agent interface {
	GetResponse(ctx context.Context, messages []Message) (string, error)
}
