//Package recipe manage the list of recipe to be applied and their workflow
package recipe

import (
	"context"
	"sort"

	"github.com/marema31/kamino/step"
	"github.com/marema31/kamino/step/common"
)

// Cooker interface for cookbook testing
type Cooker interface {
	Statistics() (map[string][]int, int)
	Load(context.Context, string, []string, []string, []string) error
	Do(context.Context) error
}

type recipe struct {
	name            string
	steps           map[uint][]common.Steper
	currentPriority uint
}

// Cookbook is a map of recipe indexed by recipe's name
type Cookbook struct {
	Recipes     map[string]recipe
	stepFactory step.Creater
}

// New returns a new
func New(sf step.Creater) *Cookbook {
	rs := make(map[string]recipe)
	return &Cookbook{
		Recipes:     rs,
		stepFactory: sf,
	}
}

// Statistics return number of step by priority by recipes and total number of steps
func (ck *Cookbook) Statistics() (map[string][]int, int) {
	result := make(map[string][]int)
	var total int
	for rname := range ck.Recipes {
		s := make([]int, 0, len(ck.Recipes[rname].steps))
		priorities := make([]int, 0, len(ck.Recipes[rname].steps))
		for priority := range ck.Recipes[rname].steps {
			priorities = append(priorities, int(priority))
		}
		sort.Ints(priorities)
		for _, priority := range priorities {
			nb := len(ck.Recipes[rname].steps[uint(priority)])
			s = append(s, nb)
			total += nb
		}
		result[rname] = s
	}
	return result, total
}

/*TODO:

Creer un package step qui sait
    - executer la step en fonction du type (fonction appellé par le moteur de recette)
     en s'appuyant sur les providers


Le package recipe doit:
	- Determiner tous les fichiers correspondant a un chaine de caractère de selection (provenant du cli)
	  en fonction des noms de repertoire/fichier
	- Creer un waitgroup, un channel en reception, Pour chaque recette de cette map, creer une goroutine qui va, dans l'ordre de priorité, a chaque niveau de priorité:
		0) recevoir le channel et l'utiliser pour l'affichage,
		1) creer un waitgroup,
		2) creer une goroutine par step de ce niveau de priorité
		3) attendre la fin du waitgroup
    - et attendre la fin du waitgroup



	Gestion des logs se baser sur Packer par exemple:
	https://github.com/hashicorp/packer/blob/3d5af49bf32aca277c573af2e454ee5ed84ef505/log.go#L17

	Ou voir a utiliser https://github.com/hashicorp/go-hclog



	Pas besoin d'utiliser un channel dans ce cas là, un mutex pour le screen est suffisant (on peut aussi utiliser channel et une goroutine dédiée)

*/
