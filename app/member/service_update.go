package member

import (
	"context"
	"fmt"
)

func (s *service) Update(ctx context.Context, username string, m Member) error {
	_, found, err := s.storage.Member(ctx, username)
	if err != nil {
		return fmt.Errorf("remove member: %w", err)
	}
	if !found {
		return ErrorMemberNotFound
	}

	m.Username = username
	err = s.storage.Update(ctx, m)
	if err != nil {
		return fmt.Errorf("update member: %w", err)
	}
	return err
}
