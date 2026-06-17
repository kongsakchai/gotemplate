package member

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/template/pkg/errs"
)

type storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *storage {
	return &storage{db: db}
}

type memberRecord struct {
	Username     string    `db:"username"`
	FirstName    string    `db:"first_name"`
	LastName     string    `db:"last_name"`
	Birthday     time.Time `db:"birthday"`
	RegisterDate time.Time `db:"register_date"`
}

func (m memberRecord) ToMember() Member {
	return Member(m)
}

func (s *storage) Members(ctx context.Context) ([]Member, error) {
	var result []memberRecord
	err := s.db.SelectContext(ctx, &result, "SELECT * FROM member")

	members := []Member{}
	for _, m := range result {
		members = append(members, m.ToMember())
	}

	return members, errs.From(err)
}

func (s *storage) Member(ctx context.Context, username string) (Member, bool, error) {
	member := memberRecord{}
	err := s.db.GetContext(ctx, &member, "SELECT * FROM member WHERE username = ?", username)
	if err == sql.ErrNoRows {
		return member.ToMember(), false, nil
	}
	return member.ToMember(), err == nil, errs.From(err)
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
