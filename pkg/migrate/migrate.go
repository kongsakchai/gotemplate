package migrate

import (
	"database/sql"

	"github.com/kongsakchai/gotemplate/pkg/config"
	mg "github.com/kongsakchai/simple-sql-migrate"
)

func Migrate(db *sql.DB, cfg config.Migration) {
	for _, v := range cfg.Versions {
		err := mg.Migrate(db, mg.Options{
			Source:  cfg.Source,
			Version: v,
		})

		if err != nil {
			panic(err)
		}
	}
}
