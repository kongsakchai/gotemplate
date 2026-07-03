package member

import (
	"context"
)

func (s *service) GetMember(ctx context.Context, username string) (Member, error) {
	members, found, err := s.storage.GetMember(ctx, username)
	if err != nil {
		return Member{}, err
	}
	if !found {
		return Member{}, ErrMemberNotFound
	}
	return members, nil
}
