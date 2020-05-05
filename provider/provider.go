package provider

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
)

//Provider provides Loader and Saver objects adapted to the datasource.
type Provider interface {
	NewLoader(context.Context, *logrus.Entry, datasource.Datasourcer, string, string) (Loader, error)
	NewSaver(context.Context, *logrus.Entry, datasource.Datasourcer, string, string, string) (Saver, error)
}

//KaminoProvider implement the Provider interface with action on database and files.
type KaminoProvider struct{}
