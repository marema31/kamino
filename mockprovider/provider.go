package mockprovider

import (
	"context"

	"github.com/marema31/kamino/datasource"
)

//MockProvider implement the Provider interface with mocked actions
type MockProvider struct {
	ErrorLoader error
	ErrorSaver  error
}

//NewLoader analyze the datasource and return mock object implementing Loader
func (p *MockProvider) NewLoader(ctx context.Context, ds datasource.Datasourcer, table string, where string) (*MockLoader, error) {
	k := MockLoader{}
	if p.ErrorLoader != nil {
		return &k, p.ErrorLoader
	}

	return &k, nil
}

//NewSaver analyze the datasource and return mock object implementing Saver
func (p *MockProvider) NewSaver(ctx context.Context, ds datasource.Datasourcer, table string, key string, mode string) (*MockSaver, error) {
	k := MockSaver{}
	if p.ErrorSaver != nil {
		return &k, p.ErrorSaver
	}

	return &k, nil
}
