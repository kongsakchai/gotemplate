package member

import (
	"context"
	"fmt"
)

func (s *service) Member(ctx context.Context, username string) (Member, error) {
	members, found, err := s.storage.Member(ctx, username)
	if err != nil {
		return Member{}, fmt.Errorf("get member by username: %w", err)
	}
	if !found {
		return Member{}, ErrorMemberNotFound
	}
	return members, nil
}
