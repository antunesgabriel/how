package agent

import (
	"context"
)

type Agent interface {
	GetResponse(ctx context.Context, input string) (string, error)
}
