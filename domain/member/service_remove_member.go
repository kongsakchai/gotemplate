package member

import (
	"context"
)

func (s *service) RemoveMember(ctx context.Context, username string) error {
	_, found, err := s.storage.GetMember(ctx, username)
	if err != nil {
		return err
	}
	if !found {
		return ErrMemberNotFound
	}

	return s.storage.RemoveMember(ctx, username)
}
