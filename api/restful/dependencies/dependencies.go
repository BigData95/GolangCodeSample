package dependencies

import (
	"context"
)

type DependencyManager interface {
	Initialize(ctx context.Context) error
	Cleanup(ctx context.Context) error
}
