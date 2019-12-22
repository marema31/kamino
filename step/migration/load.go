package migration

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Masterminds/sprig"
	"github.com/Sirupsen/logrus"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/viper"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/step/common"
)

var dialects = map[datasource.Engine]string{
	datasource.Postgres: "postgres",
	datasource.Mysql:    "mysql",
}

//PostLoad modify the loaded step values with the values provided in the map in argument
func (st *Step) PostLoad(log *logrus.Entry, superseed map[string]string) error {
	var err error

	if value, ok := superseed["migration.dir"]; ok {
		switch value {
		case "up":
			st.dir = migrate.Up
		case "down":
			st.dir = migrate.Down
		case "status":
			st.printOnly = true
		default:
			return fmt.Errorf("%s is not a correct direction for migration", value)
		}
	}

	if value, ok := superseed["migration.limit"]; ok {
		st.limit, err = strconv.Atoi(value)
	}

	// Use the noAdmin parameter only if the step configuration have noAdmin = false
	if value, ok := superseed["migration.noAdmin"]; !st.noAdmin && ok {
		st.noAdmin, err = strconv.ParseBool(value)
	}

	if value, ok := superseed["migration.noUser"]; ok {
		st.noUser, err = strconv.ParseBool(value)
	}

	return err
}

//Load data from step file using its viper representation a return priority and list of steps
func Load(ctx context.Context, log *logrus.Entry, recipePath string, name string, nameIndex int, v *viper.Viper, dss datasource.Datasourcers, force bool, dryRun bool, limitedTags []string) (priority uint, steps []common.Steper, err error) { //nolint: funlen
	steps = make([]common.Steper, 0, 1)
	priority = v.GetUint("priority")

	tags := v.GetStringSlice("tags")
	if len(tags) == 0 {
		tags = []string{""}
	}

	logStep := log.WithField("name", name).WithField("type", "migration")

	folderTmpl := v.GetString("folder")
	if folderTmpl == "" {
		logStep.Error("No migration folder provided")
		return 0, nil, fmt.Errorf("the step %s must have a folder that contains the migration files", name)
	}

	if !filepath.IsAbs(folderTmpl) {
		folderTmpl = filepath.Join(recipePath, folderTmpl)
	}

	tfolder, err := template.New("folder").Funcs(sprig.FuncMap()).Parse(folderTmpl)
	if err != nil {
		logStep.Error("Parsing the migration folder name failed:")
		logStep.Error(err)

		return 0, nil, fmt.Errorf("error parsing the folder of %s step: %v", name, err)
	}

	engines := v.GetStringSlice("engines")

	e, err := datasource.StringsToEngines(engines)
	if err != nil {
		logStep.Error(err)
		return 0, nil, err
	}

	queryTmpl := v.GetString("query")
	if queryTmpl == "" {
		logStep.Error("No SQL query provided")
		return 0, nil, fmt.Errorf("the step %s must have a query to be executed", name)
	}

	tquery, err := template.New("query").Funcs(sprig.FuncMap()).Parse(queryTmpl)
	if err != nil {
		logStep.Error("Parsing the SQL query template failed:")
		logStep.Error(err)

		return 0, nil, fmt.Errorf("error parsing the query of %s step: %v", name, err)
	}

	noUser := v.GetBool("nouser")
	noAdmin := v.GetBool("noadmin")

	tableUser := v.GetString("usertable")
	if tableUser == "" {
		tableUser = "kamino_user_migrations"
	}

	tableAdmin := v.GetString("admintable")
	if tableAdmin == "" {
		tableAdmin = "kamino_admin_migrations"
	}

	renderedQuery := bytes.NewBuffer(make([]byte, 0, 1024))
	renderedFolder := bytes.NewBuffer(make([]byte, 0, 1024))

	for index, datasource := range dss.Lookup(log, tags, limitedTags, []datasource.Type{datasource.Database}, e) {
		var step Step

		step.Name = fmt.Sprintf("%s:%d", name, nameIndex+index)
		step.datasource = datasource
		step.dryRun = dryRun

		if err := tfolder.Execute(renderedFolder, datasource.FillTmplValues()); err != nil {
			return 0, nil, err
		}

		step.folder = renderedFolder.String()
		if _, err := os.Stat(step.folder); err != nil {
			if os.IsNotExist(err) {
				log.Warningf("The migration folder %s does not exists at load phase, I will recheck at execution if it has not been created by a previous step execution", step.folder)
			}
		}

		renderedFolder.Reset()

		tmplValue := datasource.FillTmplValues()

		err = tquery.Execute(renderedQuery, tmplValue)
		if err != nil {
			logStep.Error("Rendering the query template failed")
			logStep.Error(err)

			return 0, nil, err
		}

		step.query = renderedQuery.String()
		renderedQuery.Reset()

		step.schema = tmplValue.Schema
		step.tableAdmin = tableAdmin
		step.tableUser = tableUser
		// Value for initial migration
		step.printOnly = false
		step.dir = migrate.Up
		step.limit = 0 // No limits
		step.noAdmin = noAdmin
		step.noUser = noUser

		step.dialect = dialects[datasource.GetEngine()]
		steps = append(steps, &step)
	}

	return priority, steps, nil
}
