package book

import (
	"context"
	"errors"
	"time"
)

var (
	ErrBookNotFound = errors.New("book not found")
)

type Book struct {
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Stock       int64     `json:"stock"`
	CreatedAt   time.Time `json:"createdAt"`
}

type AddBookPayload struct {
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"createdAt"`
}

type UpdateBookPayload struct {
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Stock       int64     `json:"stock"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Storager interface {
	GetBook(ctx context.Context, code string) (Book, bool, error)
	GetBooksByName(ctx context.Context, name string) ([]Book, error)
	AddBook(ctx context.Context, payload AddBookPayload) error
	UpdateBook(ctx context.Context, payload UpdateBookPayload) error
	UpdateStock(ctx context.Context, code string, stock int64) error
}

type Servicer interface {
	BooksByName(ctx context.Context, name string) ([]Book, error)
	AddBook(ctx context.Context, payload AddBookPayload) error
}
