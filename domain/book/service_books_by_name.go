package book

import (
	"context"
)

func (s *service) BooksByName(ctx context.Context, name string) ([]Book, error) {
	return s.storage.GetBooksByName(ctx, name)
}
