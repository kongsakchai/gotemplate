package migrate

import (
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kongsakchai/gotemplate/pkg/config"
)

func Migrate(cfg config.Migration) {
	if !cfg.Enable {
		return
	}

	m, err := migrate.New(cfg.Source, cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}
	if cfg.Version >= 0 {
		err = m.Migrate(uint(cfg.Version))
	} else {
		err = m.Up()
	}
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}

	v, dirty, err := m.Version()
	if err != nil {
		panic(err)
	}
	slog.Info("Migration", "version", v, "dirty", dirty)
}
