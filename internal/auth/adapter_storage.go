package auth

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/pkg/serror"
)

type storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *storage {
	return &storage{db: db}
}

type userRecord struct {
	Id        string    `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password_hash"`
	CreatedAt time.Time `db:"created_at"`
}

func (r userRecord) toUser() User {
	return User{
		Id:        r.Id,
		Username:  r.Username,
		Password:  r.Password,
		CreatedAt: r.CreatedAt,
	}
}

func (s *storage) FindUserByUsername(ctx context.Context, username string) (User, bool, error) {
	var user userRecord
	err := s.db.GetContext(ctx, &user, "SELECT id, username, password_hash, created_at FROM user WHERE username = ?", username)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, false, nil
		}
		return User{}, false, serror.From(err)
	}
	return user.toUser(), true, nil
}

func (s *storage) CreateUser(ctx context.Context, user User) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO user (id, username, password_hash) VALUES (?, ?, ?)", user.Id, user.Username, user.Password)
	if err != nil {
		return serror.From(err)
	}
	return nil
}
