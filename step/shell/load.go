package shell

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/step/common"
)

//PostLoad modify the loaded step values with the values provided in the map in argument.
func (st *Step) PostLoad(log *logrus.Entry, superseed map[string]string) error {
	return nil
}

//Load data from step file using its viper representation a return priority and list of steps.
func Load(ctx context.Context, log *logrus.Entry, recipePath string, name string, nameIndex int, v *viper.Viper, dss datasource.Datasourcers, force bool, dryRun bool, limitedTags []string) (priority uint, steps []common.Steper, err error) { //nolint: funlen
	steps = make([]common.Steper, 0, 1)
	priority = v.GetUint("priority")
	ignoreErrors := v.GetBool("ignoreErrors")

	tags := v.GetStringSlice("tags")
	if len(tags) == 0 {
		tags = []string{""}
	}

	logStep := log.WithField("name", name).WithField("type", "shell")

	script := v.GetString("script")
	if script == "" {
		logStep.Error("No script provided")
		return 0, nil, fmt.Errorf("the step %s must have a script to call: %w", name, common.ErrMissingParameter)
	}

	args := v.GetStringSlice("arguments")

	targs := make([]*template.Template, len(args))
	for i, arg := range args {
		targs[i], err = template.New(fmt.Sprintf("env%d", i)).Parse(arg)
		if err != nil {
			logStep.Errorf("Parsing the %dth environment failed:", i)
			logStep.Error(err)

			return 0, nil, fmt.Errorf("error parsing the environment of %s step: %w", name, err)
		}
	}

	rendered := bytes.NewBuffer(make([]byte, 0, 1024))

	path := v.GetString("path")
	if path == "" {
		path = recipePath
	} else if !filepath.IsAbs(path) {
		path = filepath.Join(recipePath, path)
	}

	tpath, err := template.New("path").Parse(path)
	if err != nil {
		logStep.Error("Parsing the path failed:")
		logStep.Error(err)

		return 0, nil, fmt.Errorf("error parsing the path of %s step: %w", name, err)
	}

	envs := v.GetStringSlice("environment")
	tenvs := make([]*template.Template, len(envs))

	for i, env := range envs {
		tenvs[i], err = template.New(fmt.Sprintf("env%d", i)).Parse(env)
		if err != nil {
			logStep.Errorf("Parsing the %dth environment failed:", i)
			logStep.Error(err)

			return 0, nil, fmt.Errorf("error parsing the environment of %s step: %w", name, err)
		}
	}

	engines := v.GetStringSlice("engines")

	e, err := datasource.StringsToEngines(engines)
	if err != nil {
		logStep.Error(err)
		return 0, nil, err
	}

	lookedUp, _ := dss.Lookup(log, tags, limitedTags, []datasource.Type{datasource.Database}, e)
	for index, datasource := range lookedUp {
		var step Step

		step.Name = fmt.Sprintf("%s:%d", name, nameIndex+index)
		step.dryRun = dryRun
		step.ignoreErrors = ignoreErrors
		cmdPath := script

		tmplValues := datasource.FillTmplValues()

		if err := tpath.Execute(rendered, tmplValues); err != nil {
			return 0, nil, err
		}

		cmdDir := rendered.String()
		rendered.Reset()

		args = make([]string, 1)
		args[0] = cmdPath

		for _, targ := range targs {
			if err := targ.Execute(rendered, tmplValues); err != nil {
				return 0, nil, err
			}

			args = append(args, rendered.String())
			rendered.Reset()
		}

		cmdArgs := args

		envs = make([]string, 0)

		for _, tenv := range tenvs {
			if err := tenv.Execute(rendered, tmplValues); err != nil {
				return 0, nil, err
			}

			envs = append(envs, rendered.String())
			rendered.Reset()
		}

		cmdEnv := envs

		step.cmd = &exec.Cmd{
			Path: cmdPath,
			Dir:  cmdDir,
			Args: cmdArgs,
			Env:  cmdEnv,
		}

		realScript := filepath.Join(cmdDir, cmdPath)
		if _, err := os.Stat(realScript); err != nil {
			if os.IsNotExist(err) {
				log.Warningf("The script %s does not exists at load phase, I will recheck at execution if it has not been created by a previous step execution", realScript)
			}
		}

		step.datasource = datasource

		steps = append(steps, &step)
	}

	return priority, steps, nil
}
