package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
)

func isRecipeFolder(log *logrus.Entry, filename string) bool {
	complete := true
	file, err := os.Stat(path.Join(CfgFolder, filename))

	if os.IsNotExist(err) {
		log.Infof("%s does not %s recipe", CfgFolder, filename)
		return false
	}

	if !file.Mode().IsDir() {
		log.Infof("%s is not a recipe folder", filename)
		return false
	}

	for _, subFolder := range []string{"datasources", "steps"} {
		info, err := os.Stat(path.Join(CfgFolder, filename, subFolder))
		if os.IsNotExist(err) {
			log.Infof("%s does not contains %s", filename, subFolder)

			complete = false
		} else if !info.IsDir() {
			log.Infof("%s for %s is not a folder", subFolder, filename)

			complete = false
		}
	}

	return complete
}

func lookupProvidedRecipes(log *logrus.Entry, recipesNames []string) ([]string, error) {
	recipes := make([]string, 0, len(recipesNames))

	for _, filename := range recipesNames {
		if isRecipeFolder(log, filename) {
			recipes = append(recipes, filename)
			continue // We found, go to the next arg
		}

		if !strings.EqualFold(filepath.Base(filename), filename) {
			return nil, fmt.Errorf("%s (%s) is not a valid globbed recipe: %w", filename, filepath.Base(filename), ErrWrongParameterValue)
		}

		folders, err := filepath.Glob(filepath.Join(CfgFolder, filename))

		if err != nil {
			return nil, fmt.Errorf("%s is not a valid globbed recipe: %w", filename, err)
		}

		found := false

		for _, recipe := range folders {
			recipe = "." + strings.TrimPrefix(recipe, CfgFolder)
			recipe = filepath.Base(recipe)

			if isRecipeFolder(log, recipe) {
				recipes = append(recipes, recipe)
				found = true
			}
		}

		if !found {
			return nil, fmt.Errorf("no recipe correspond to %s: %w", filename, ErrWrongParameterValue)
		}
	}

	return recipes, nil
}

func lookupAllRecipes(log *logrus.Entry) ([]string, error) {
	recipes := make([]string, 0)

	files, err := ioutil.ReadDir(CfgFolder)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.Mode().IsDir() && isRecipeFolder(log, file.Name()) {
			recipes = append(recipes, file.Name())
		}
	}

	if len(recipes) == 0 {
		return nil, fmt.Errorf("no recipes folder found in %s: %w", CfgFolder, ErrNoConfiguration)
	}

	return recipes, nil
}

// FindRecipes lookup the configuration folder and return a list of recipes if the args is empty.
func FindRecipes(log *logrus.Entry, args []string) ([]string, error) {
	if len(args) != 0 {
		return lookupProvidedRecipes(log, args)
	}

	return lookupAllRecipes(log)
}
