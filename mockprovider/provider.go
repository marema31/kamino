package mockprovider

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider"
)

//MockProvider implement the Provider interface with mocked actions
type MockProvider struct {
	ErrorLoader error
	ErrorSaver  error
	Loader      *MockLoader
	Savers      []*MockSaver
}

//NewLoader analyze the datasource and return mock object implementing Loader
func (p *MockProvider) NewLoader(ctx context.Context, log *logrus.Entry, ds datasource.Datasourcer, table string, where string) (provider.Loader, error) {
	if p.ErrorLoader != nil {
		return nil, p.ErrorLoader
	}
	k := &MockLoader{}
	p.Loader = k
	return k, nil
}

//NewSaver analyze the datasource and return mock object implementing Saver
func (p *MockProvider) NewSaver(ctx context.Context, log *logrus.Entry, ds datasource.Datasourcer, table string, key string, mode string) (provider.Saver, error) {
	if p.ErrorSaver != nil {
		return nil, p.ErrorSaver
	}
	k := MockSaver{}
	p.Savers = append(p.Savers, &k)
	return &k, nil
}
