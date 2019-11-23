package common

import (
	"fmt"
	"io/ioutil"
)

// FindRecipes lookup the configuration folder and return a list of recipes if the args is empty
func FindRecipes(args []string) ([]string, error) {
	if len(args) != 0 {
		return args, nil
	}
	recipes := make([]string, 0)
	files, err := ioutil.ReadDir(CfgFolder)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.Mode().IsDir() {
			recipes = append(recipes, file.Name())
		}
	}
	if len(recipes) == 0 {
		return nil, fmt.Errorf("no recipes folder found in %s", CfgFolder)
	}
	return recipes, nil
}
