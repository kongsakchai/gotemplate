package sqldb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "modernc.org/sqlite"
)

const errPingDriverName = "sqldb_test_err_ping"

func init() {
	sql.Register(errPingDriverName, &errPingDriver{})
}

// errPingConn is a driver.Conn that implements driver.Pinger and always fails to ping.
type errPingConn struct{}

func (errPingConn) Ping(ctx context.Context) error {
	return errors.New("ping fail")
}

func (errPingConn) Prepare(query string) (driver.Stmt, error) {
	return nil, errors.New("not implemented")
}

func (errPingConn) Close() error {
	return nil
}

func (errPingConn) Begin() (driver.Tx, error) {
	return nil, errors.New("not implemented")
}

// errPingDriver is a driver.Driver whose Open succeeds but connection ping always fails.
type errPingDriver struct{}

func (errPingDriver) Open(name string) (driver.Conn, error) {
	return errPingConn{}, nil
}

func TestNewDatabase(t *testing.T) {
	t.Run("should panic when unknow driver", func(t *testing.T) {
		// check panic
		defer func() {
			p := recover()
			assert.NotNil(t, p)
			assert.Contains(t, p.(string), "Connect to database error")
		}()

		New("unknow", "invalid")
	})

	t.Run("should panic when ping database error", func(t *testing.T) {
		defer func() {
			p := recover()
			assert.NotNil(t, p)
			assert.Contains(t, p.(string), "Ping database error")
		}()

		New(errPingDriverName, "ignored")
	})

	t.Run("should ping success when db connet", func(t *testing.T) {
		// check panic
		defer func() {
			p := recover()
			assert.Nil(t, p)
		}()

		db, close := New("sqlite", ":memory:")
		defer close(t.Context())

		err := db.Ping()

		assert.NoError(t, err)
	})
}
