package migration

import (
	"fmt"
	"path"

	"github.com/Sirupsen/logrus"
	migrate "github.com/rubenv/sql-migrate"
)

func (st *Step) apply(log *logrus.Entry, admin bool, limit int) (int, error) {
	folder := st.folder
	tableName := st.tableUser

	if admin {
		folder = path.Join(st.folder, "admin")
		tableName = st.tableAdmin
	}

	direction := "Up"
	if st.dir == migrate.Down {
		direction = "Down"
	}

	db, err := st.datasource.OpenDatabase(log, admin, false)
	if err != nil {
		return 0, err
	}

	migSet := migrate.MigrationSet{TableName: tableName, SchemaName: st.schema}
	source := migrate.FileMigrationSource{Dir: folder}

	if st.dryRun {
		migrations, _, err := migSet.PlanMigration(db, st.dialect, source, st.dir, limit)
		if err != nil {
			log.Error("Planning migration failed")
			log.Error(err)

			return 0, err
		}

		for _, m := range migrations {
			log.Infof("Would apply migration %s (%s)", m.Id, direction)
			statements := m.Up

			if st.dir == migrate.Down {
				statements = m.Down
			}

			for _, q := range statements {
				log.Info(q)
			}
		}

		return 0, nil
	}

	applied, err := migSet.ExecMax(db, st.dialect, source, st.dir, limit)
	if err != nil {
		log.Error("Applying migration failed")
		log.Error(err)

		return applied, err
	}

	log.Infof("Applied %d %s migration(s)", applied, direction)

	migrationsLeft, _, err := migSet.PlanMigration(db, st.dialect, source, st.dir, limit)
	if err != nil {
		log.Error("Planning migration failed")
		log.Error(err)

		return applied, err
	}

	if limit == 0 && len(migrationsLeft) > 0 {
		log.Errorf("All migrations should have applied but only %d was applied (%d not applied)", applied, len(migrationsLeft))
		return applied, fmt.Errorf("all migrations should have applied")
	}

	if len(migrationsLeft) > 0 && limit > applied {
		log.Errorf("%d migrations should have applied but only %d was applied (up to %d could be applied)", limit, applied, len(migrationsLeft))
		return applied, fmt.Errorf("%d migrations should have applied", limit)
	}

	return applied, nil
}
