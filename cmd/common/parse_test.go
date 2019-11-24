package common_test

import (
	"testing"

	"github.com/marema31/kamino/cmd/common"
)

func TestFindRecipeArgs(t *testing.T) {
	recipes, err := common.FindRecipes([]string{"r1", "r2"})
	if err != nil {
		t.Fatalf("Should not return error, returned %v", err)
	}
	if len(recipes) == 0 {
		t.Fatalf("No recipes found")
	}
	if len(recipes) != 2 {
		t.Fatalf("Waiting for 'r1','r2' returned %v", recipes)
	}
}

func TestFindRecipeFolderOk(t *testing.T) {
	common.CfgFolder = "testdata/good"
	recipes, err := common.FindRecipes([]string{})
	if err != nil {
		t.Fatalf("Should not return error, returned %v", err)
	}
	if len(recipes) == 0 {
		t.Fatalf("No recipes found")
	}
	if len(recipes) != 3 {
		t.Fatalf("Waiting for 'rec1','rec2', 'rec3' returned %v", recipes)
	}
}

func TestFindRecipeFolderEmpty(t *testing.T) {
	common.CfgFolder = "testdata/fail"
	recipes, err := common.FindRecipes([]string{})
	if err == nil {
		t.Fatalf("Should return error")
	}
	if len(recipes) != 0 {
		t.Fatalf("Recipes %v found, none should have been found", recipes)
	}
}

func TestFindRecipeFolderError(t *testing.T) {
	common.CfgFolder = "testdata/notexist"
	recipes, err := common.FindRecipes([]string{})
	if err == nil {
		t.Fatalf("Should return error")
	}
	if len(recipes) != 0 {
		t.Fatalf("Recipes %v found, none should have been found", recipes)
	}
}
