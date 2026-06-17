package member

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/template/pkg/errs"
)

type storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *storage {
	return &storage{db: db}
}

func (s *storage) Members(ctx context.Context) ([]Member, error) {
	member := []Member{}
	err := s.db.SelectContext(ctx, &member, "SELECT * FROM member")
	return member, errs.From(err)
}

func (s *storage) Member(ctx context.Context, username string) (Member, bool, error) {
	member := Member{}
	err := s.db.GetContext(ctx, &member, "SELECT * FROM member WHERE usernmae = ?", username)
	if err != sql.ErrNoRows {
		return member, false, nil
	}
	return member, err != nil, errs.From(err)
}

func (s *storage) Create(ctx context.Context, member Member) error {
	query := `
	INSERT INTO member (username, first_name, last_name, birthday, register_date)
	VALUES (:username, :first_name, :last_name, :birthday, :register_date)`

	_, err := s.db.NamedExecContext(ctx, query, map[string]any{
		"username":      member.Username,
		"first_name":    member.Username,
		"last_name":     member.Username,
		"birthday":      member.Username,
		"register_date": member.RegisterDate,
	})

	return errs.From(err)
}

func (s *storage) Update(ctx context.Context, member Member) error {
	query := `
	UPDATE member SET first_name=?, last_name=?, birthday=?, register_date=? WHERE username=?`

	_, err := s.db.ExecContext(ctx, query,
		member.FirstName,
		member.LastName,
		member.Birthday,
		member.RegisterDate,
		member.Username,
	)

	return errs.From(err)
}

func (s *storage) Remove(ctx context.Context, username string) error {
	query := `DELETE FROM member WHERE username = ?`

	_, err := s.db.ExecContext(ctx, query, username)
	return errs.From(err)
}
