//Package common provides the utility functions and type needed by all specialized step packages
package common

import (
	"context"

	"github.com/marema31/kamino/datasource"
)

// Steper Interface that will be used to run the steps
type Steper interface {
	Do(context.Context) error
}

// TmplValues structure use for template rendering to avoid exposing the datasource structure to the template
type TmplValues struct {
	Database string
	User     string
	Password string
	Schema   string
	Host     string
	Port     string
	Tags     []string
	//TODO: add named tags
}

// FillTmplValues return a struct for template operation with value corresponding to the provided datasource
func FillTmplValues(ds *datasource.Datasource) TmplValues {
	var tv TmplValues

	tv.Database = ds.Database
	tv.User = ds.User
	tv.Password = ds.UserPw
	tv.Schema = ds.Schema
	tv.Host = ds.Host
	tv.Port = ds.Port
	tv.Tags = ds.Tags

	return tv
}
