package mockdatasource

import (
	"fmt"
	"strings"

	"github.com/marema31/kamino/datasource"
)

// MockDatasources is fake datasources object for test purpose
type MockDatasources struct {
	// Datasource tag dictionnary for lookup
	precordedAnswers map[string][]*MockDatasource
}

// New returns a new Datasources object with elments initialized
func New() *MockDatasources {
	var dss MockDatasources

	dss.precordedAnswers = make(map[string][]*MockDatasource)
	return &dss
}

//LoadAll do nothing, return an error if path is empty
func (dss *MockDatasources) LoadAll(path string) error {
	if path == "" {
		return fmt.Errorf("empty path")
	}
	return nil
}

func getIndex(tags []string, dsTypes []datasource.Type, engines []datasource.Engine) string {
	index := strings.Join(tags, "@")
	index = index + "&"
	for _, t := range dsTypes {
		index = index + "@" + datasource.TypeToString(t)
	}
	index = index + "&"
	for _, t := range engines {
		index = index + "@" + datasource.EngineToString(t)
	}
	return index
}

//Lookup : return the corresponding array of Mocked datasources
//WARNING: the algorithm of lookup is much simpler than the one from the object, all the parameters must be exactly the same !
func (dss *MockDatasources) Lookup(tags []string, dsTypes []datasource.Type, engines []datasource.Engine) []datasource.Datasourcer {
	index := getIndex(tags, dsTypes, engines)
	dsr := make([]datasource.Datasourcer, 0, len(dss.precordedAnswers[index]))
	for _, ds := range dss.precordedAnswers[index] {
		dsr = append(dsr, ds)
	}
	return dsr
}

//Insert add a Mocked datasource to the array
func (dss *MockDatasources) Insert(tags []string, dsTypes []datasource.Type, engines []datasource.Engine, ds []*MockDatasource) {
	index := getIndex(tags, dsTypes, engines)
	dss.precordedAnswers[index] = ds
}
