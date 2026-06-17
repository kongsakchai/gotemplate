package member

import (
	"context"
	"fmt"
)

func (s *service) Members(ctx context.Context) ([]Member, error) {
	members, err := s.storage.Members(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all member: %w", err)
	}
	return members, nil
}
