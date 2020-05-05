package mockprovider

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider"
)

//MockProvider implement the Provider interface with mocked actions.
type MockProvider struct {
	ErrorLoader   error
	LoaderToFail  int
	CurrentLoader int
	ErrorSaver    error
	SaverToFail   int
	CurrentSaver  int
	Loader        *MockLoader
	Savers        []*MockSaver
}

//NewLoader analyze the datasource and return mock object implementing Loader.
func (p *MockProvider) NewLoader(ctx context.Context, log *logrus.Entry, ds datasource.Datasourcer, table string, where string) (provider.Loader, error) {
	if p.ErrorLoader != nil && p.CurrentLoader == p.LoaderToFail {
		err := p.ErrorLoader
		p.CurrentLoader++

		return nil, err
	}

	k := &MockLoader{}
	p.Loader = k
	p.CurrentLoader++

	return k, nil
}

//NewSaver analyze the datasource and return mock object implementing Saver.
func (p *MockProvider) NewSaver(ctx context.Context, log *logrus.Entry, ds datasource.Datasourcer, table string, key string, mode string) (provider.Saver, error) {
	if p.ErrorSaver != nil && p.CurrentSaver == p.SaverToFail {
		p.CurrentSaver++
		return nil, p.ErrorSaver
	}

	k := MockSaver{}
	p.Savers = append(p.Savers, &k)
	p.CurrentSaver++

	return &k, nil
}
