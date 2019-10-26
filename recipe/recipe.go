//Package recipe manage the list of recipe to be applied and their workflow
package recipe

/*TODO:

Creer un package step qui sait
   - load un fichier donné par le moteur de recette (en s'appuyant sur le datasource.Lookup)
	 et renvoyer une priorité et
		une liste de step concurrente  (une par datasource pour sql/migration, une seule par step pour sync)
		la selection des datasource se fait à cause des engines
			et tags de ce step filtré par la liste des tags passés par le moteur de recette

   - executer la step en fonction du type (fonction appellé par le moteur de recette)
     en s'appuyant sur les providers


Le package recipe doit:
	- Determiner tous les fichiers correspondant a un chaine de caractère de selection (provenant du cli)
	  en fonction des noms de repertoire/fichier
	- Creer un map[recette]map[priorité] de slice de *step en chargeant toutes les steps correspondantes

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
