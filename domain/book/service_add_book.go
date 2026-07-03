package book

import (
	"context"
)

func (s *service) Add(ctx context.Context, payload AddBookPayload) error {
	book, found, err := s.storage.GetBook(ctx, payload.Code)
	if err != nil {
		return err
	}
	if found {
		return s.storage.UpdateStock(ctx, book.Code, book.Stock+1)
	}

	return s.storage.AddBook(ctx, payload)
}
