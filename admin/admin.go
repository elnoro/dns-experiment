package admin

import (
	"context"
)

type DbServer interface {
	Run(ctx context.Context) error
}
