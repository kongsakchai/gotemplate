package member

import (
	"context"
	"fmt"
	"time"
)

func (s *service) Create(ctx context.Context, m Member) error {
	m.RegisterDate = s.clock.Now()

	age := m.RegisterDate.Sub(m.Birthday)
	if age < 15*365*24*time.Hour {
		return ErrorMinAge
	}
	if age > 60*365*24*time.Hour {
		return ErrorMaxAge
	}

	_, exiting, err := s.storage.Member(ctx, m.Username)
	if err != nil {
		return fmt.Errorf("create member: %w", err)
	}
	if exiting {
		return ErrorDuplicate
	}

	return s.storage.Create(ctx, m)
}
