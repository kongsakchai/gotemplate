package todo

import (
	"github.com/kongsakchai/gotemplate/app"
	"github.com/stretchr/testify/mock"
)

type mockService struct {
	mock.Mock
	Servicer
}

func (m *mockService) Todos(ctx app.Context) ([]Todo, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Todo), args.Error(1)
}

func (m *mockService) Todo(ctx app.Context, id string) (Todo, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Todo), args.Error(1)
}

func (m *mockService) Create(ctx app.Context, todo Todo) (Todo, error) {
	args := m.Called(ctx, todo)
	return args.Get(0).(Todo), args.Error(1)
}

func (m *mockService) Update(ctx app.Context, todo Todo) (Todo, error) {
	args := m.Called(ctx, todo)
	return args.Get(0).(Todo), args.Error(1)
}

func (m *mockService) Delete(ctx app.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
