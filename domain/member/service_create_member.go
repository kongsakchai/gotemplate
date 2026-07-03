package member

import (
	"context"
	"time"
)

func (s *service) CreateMember(ctx context.Context, m CreateMemberPayload) error {
	age := s.clock.Now().Sub(m.Birthday)
	if age < 15*365*24*time.Hour || age > 60*365*24*time.Hour {
		return ErrInvalidAge.Causef("birth at %s age", m.Birthday.Format("2006-01-02"))
	}

	_, exiting, err := s.storage.GetMember(ctx, m.Username)
	if err != nil {
		return err
	}
	if exiting {
		return ErrUsernameUnavailable
	}

	return s.storage.CreateMember(ctx, m)
}
