package member

import (
	"context"
)

func (s *service) UpdateMember(ctx context.Context, m UpdateMemberPayload) error {
	_, found, err := s.storage.GetMember(ctx, m.Username)
	if err != nil {
		return err
	}
	if !found {
		return ErrMemberNotFound
	}

	err = s.storage.UpdateMember(ctx, m)
	if err != nil {
		return err
	}
	return err
}
