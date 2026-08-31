package todo

import "context"

type service struct {
	st Storager
}

func NewService(st Storager) *service {
	return &service{st: st}
}

func (s *service) CreateTodo(ctx context.Context, userId string, todo Todo) error {
	hasUser, err := s.st.HasUser(ctx, userId)
	if err != nil {
		return err
	}
	if !hasUser {
		return ErrUserNotFound.Err()
	}
	todo.UserID = userId
	return s.st.CreateTodo(ctx, todo)
}

func (s *service) GetTodos(ctx context.Context, userID string) ([]Todo, error) {
	return s.st.GetTodos(ctx, userID)
}

func (s *service) UpdateTodo(ctx context.Context, userId string, todo Todo) error {
	_, found, err := s.st.FindTodo(ctx, userId, todo.ID)
	if err != nil {
		return err
	}
	if !found {
		return ErrTodoNotFound.Err()
	}
	todo.UserID = userId
	return s.st.UpdateTodo(ctx, todo)
}

func (s *service) DeleteTodo(ctx context.Context, userId string, id string) error {
	_, found, err := s.st.FindTodo(ctx, userId, id)
	if err != nil {
		return err
	}
	if !found {
		return ErrTodoNotFound.Err()
	}
	return s.st.DeleteTodo(ctx, id)
}
