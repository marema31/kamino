package mockdatasource

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
)

// MockDatasources is fake datasources object for test purpose.
type MockDatasources struct {
	// Datasource tag dictionary for lookup
	precordedAnswersLimited    map[string][]*MockDatasource
	precordedAnswersNotLimited map[string][]*MockDatasource
}

// New returns a new Datasources object with elments initialized.
func New() *MockDatasources {
	var dss MockDatasources

	dss.precordedAnswersLimited = make(map[string][]*MockDatasource)
	dss.precordedAnswersNotLimited = make(map[string][]*MockDatasource)

	return &dss
}

var errNoConfiguration = errors.New("NO CONFIGURATION FOUND")

//LoadAll do nothing, return an error if path is empty.
func (dss *MockDatasources) LoadAll(path string, log *logrus.Entry) error {
	if path == "" {
		return fmt.Errorf("empty path: %w", errNoConfiguration)
	}

	return nil
}

//CloseAll do nothing, return an error if path is empty.
func (dss *MockDatasources) CloseAll(log *logrus.Entry) {
}

func getIndex(tags []string, dsTypes []datasource.Type, engines []datasource.Engine) string {
	index := strings.Join(tags, "@")
	index += "&"

	for _, t := range dsTypes {
		index = index + "@" + datasource.TypeToString(t)
	}

	index += "&"

	for _, t := range engines {
		index = index + "@" + datasource.EngineToString(t)
	}

	return index
}

//Lookup : return the corresponding array of Mocked datasources
//WARNING: the algorithm of lookup is much simpler than the one from the object, all the parameters must be exactly the same !
func (dss *MockDatasources) Lookup(log *logrus.Entry, tags []string, limitedTags []string, dsTypes []datasource.Type, engines []datasource.Engine) ([]datasource.Datasourcer, []datasource.Datasourcer) {
	index := getIndex(tags, dsTypes, engines)
	dsrl := make([]datasource.Datasourcer, 0, len(dss.precordedAnswersLimited[index]))

	for _, ds := range dss.precordedAnswersLimited[index] {
		dsrl = append(dsrl, ds)
	}

	dsrn := make([]datasource.Datasourcer, 0, len(dss.precordedAnswersNotLimited[index]))

	for _, ds := range dss.precordedAnswersNotLimited[index] {
		dsrn = append(dsrn, ds)
	}

	return dsrl, dsrn
}

//Insert add a Mocked datasource to the array.
func (dss *MockDatasources) Insert(limited bool, tags []string, dsTypes []datasource.Type, engines []datasource.Engine, ds []*MockDatasource) {
	index := getIndex(tags, dsTypes, engines)

	if limited {
		dss.precordedAnswersLimited[index] = ds
		return
	}

	dss.precordedAnswersNotLimited[index] = ds
}
