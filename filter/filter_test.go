package filter_test

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/filter"
	"github.com/marema31/kamino/provider/types"
)

func TestFilterUnknown(t *testing.T) {
	aParams := []string{}
	mParams := make(map[string]string)
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)

	_, err := filter.NewFilter(log, "unknown", aParams, mParams)
	if err == nil {
		t.Errorf("NewFilter should returns an error")
	}
}

func TestFilterOnlyOk(t *testing.T) {
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)
	aParams := []string{"id", "name", "sex"}
	mParams := make(map[string]string)
	f, err := filter.NewFilter(log, "only", aParams, mParams)
	if err != nil {
		t.Errorf("NewFilter should not returns an error, returned: %v", err)
	}

	in := make(types.Record)
	in["id"] = "1"
	in["name"] = "Doe"
	in["firstname"] = "John"
	out, err := f.Filter(in)
	if err != nil {
		t.Errorf("Filter should not returns an error, returned: %v", err)
	}

	id, ok := out["id"]
	if !ok {
		t.Errorf("The filtered result should contain 'id' columns")
	}
	if id != "1" {
		t.Errorf("The filtered name should be '1', it is '%s'", id)
	}

	name, ok := out["name"]
	if !ok {
		t.Errorf("The filtered result should contain 'name' columns")
	}
	if name != "Doe" {
		t.Errorf("The filtered name should be 'Doe', it is '%s'", name)
	}

	if _, ok := out["firstname"]; ok {
		t.Errorf("The filtered result should not contain 'firstname' columns")
	}

	if _, ok := out["sex"]; ok {
		t.Errorf("The filtered result should not contain 'sex' columns")
	}
}

func TestFilterOnlyFail(t *testing.T) {
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)
	aParams := []string{}
	mParams := make(map[string]string)
	_, err := filter.NewFilter(log, "only", nil, mParams)
	if err == nil {
		t.Errorf("NewFilter only without parameters should returns an error")
	}
	_, err = filter.NewFilter(log, "only", aParams, mParams)
	if err == nil {
		t.Errorf("NewFilter only without parameters should returns an error")
	}
}

func TestFilterReplaceOk(t *testing.T) {
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)
	aParams := []string{}
	mParams := make(map[string]string)
	mParams["id"] = "42"
	mParams["firstname"] = "Jane"
	mParams["sex"] = "female"

	f, err := filter.NewFilter(log, "replace", aParams, mParams)
	if err != nil {
		t.Errorf("NewFilter should not returns an error, returned: %v", err)
	}

	in := make(types.Record)
	in["id"] = "1"
	in["name"] = "Doe"
	in["firstname"] = "John"
	out, err := f.Filter(in)
	if err != nil {
		t.Errorf("Filter should not returns an error, returned: %v", err)
	}

	id, ok := out["id"]
	if !ok {
		t.Errorf("The filtered result should contain 'id' columns")
	}
	if id != "42" {
		t.Errorf("The filtered name should be '42', it is '%s'", id)
	}

	name, ok := out["name"]
	if !ok {
		t.Errorf("The filtered result should contain 'name' columns")
	}
	if name != "Doe" {
		t.Errorf("The filtered name should be 'Doe', it is '%s'", name)
	}

	firstname, ok := out["firstname"]
	if !ok {
		t.Errorf("The filtered result should contain 'firstname' columns")
	}
	if firstname != "Jane" {
		t.Errorf("The filtered firstname should be 'Jane', it is '%s'", firstname)
	}

	sex, ok := out["sex"]
	if !ok {
		t.Errorf("The filtered result should contain 'sex' columns")
	}
	if sex != "female" {
		t.Errorf("The filtered sex should be 'female', it is '%s'", sex)
	}
}

func TestFilterReplaceFail(t *testing.T) {
	logger := logrus.New()
	log := logger.WithField("appname", "kamino")
	logger.SetLevel(logrus.PanicLevel)
	aParams := []string{}
	mParams := make(map[string]string)
	_, err := filter.NewFilter(log, "replace", aParams, nil)
	if err == nil {
		t.Errorf("NewFilter replace without parameters should returns an error")
	}
	_, err = filter.NewFilter(log, "replace", aParams, mParams)
	if err == nil {
		t.Errorf("NewFilter replace without parameters should returns an error")
	}
}
