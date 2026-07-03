package member

import (
	"context"
)

func (s *service) GetMembers(ctx context.Context) ([]Member, error) {
	return s.storage.GetMembers(ctx)
}
