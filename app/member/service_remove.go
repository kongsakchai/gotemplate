package member

import (
	"context"
	"fmt"
)

func (s *service) Remove(ctx context.Context, username string) error {
	_, found, err := s.storage.Member(ctx, username)
	if err != nil {
		return fmt.Errorf("remove member: %w", err)
	}
	if !found {
		return ErrorMemberNotFound
	}

	err = s.storage.Remove(ctx, username)
	if err != nil {
		return fmt.Errorf("remove member: %w", err)
	}
	return err
}
