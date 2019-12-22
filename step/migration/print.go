package migration

import (
	"fmt"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/olekukonko/tablewriter"
	migrate "github.com/rubenv/sql-migrate"
)

func (st *Step) print(log *logrus.Entry, admin bool) (int, error) {
	folder := st.folder
	tableName := st.tableUser
	caption := fmt.Sprintf("%s (%s)", st.datasource.GetName(), "user")

	if admin {
		folder = path.Join(st.folder, "admin")
		caption = fmt.Sprintf("%s (%s)", st.datasource.GetName(), "admin")
		tableName = st.tableAdmin
	}

	db, err := st.datasource.OpenDatabase(log, admin, false)
	if err != nil {
		return 0, err
	}

	migSet := migrate.MigrationSet{TableName: tableName, SchemaName: st.schema}
	source := migrate.FileMigrationSource{Dir: folder}

	files, err := source.FindMigrations()
	if err != nil {
		log.Error(err)
		return 0, err
	}

	records, err := migSet.GetMigrationRecords(db, st.dialect)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Migration", "Applied"})
	table.SetColWidth(60)

	rows := make(map[string]string)

	// Initialize table with migrations files
	for _, f := range files {
		rows[f.Id] = "no"
	}

	//Find the applying date
	for _, r := range records {
		if _, present := rows[r.Id]; !present {
			log.Errorf("Applied migration %s not exists", r.Id)
		}

		rows[r.Id] = r.AppliedAt.String()
	}

	remaining := 0

	for id, applied := range rows {
		table.Append([]string{
			id,
			applied,
		})

		if applied == "no" {
			remaining++
		}
	}

	fmt.Printf("\n  ---------- %s ----------\n", caption)
	table.Render()
	fmt.Println()

	return remaining, nil
}
