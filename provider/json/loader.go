package json

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/marema31/kamino/datasource"
	"github.com/marema31/kamino/provider/types"
)

//KaminoJSONLoader specifc state for database Saver provider
type KaminoJSONLoader struct {
	ds         datasource.Datasourcer
	file       io.ReadCloser
	name       string
	content    []map[string]string
	currentRow int
}

//NewLoader open the encoding process on provider file, read the whole file, unMarshal it and return a Loader compatible object
func NewLoader(ctx context.Context, ds datasource.Datasourcer) (*KaminoJSONLoader, error) {
	file, err := ds.OpenReadFile()
	if err != nil {
		return nil, err
	}
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	tv := ds.FillTmplValues()
	k := KaminoJSONLoader{ds: ds, file: file, name: tv.FilePath, content: nil, currentRow: 0}
	k.content = make([]map[string]string, 0)
	err = json.Unmarshal([]byte(byteValue), &k.content)
	if err != nil {
		return nil, err
	}

	return &k, nil
}

//Next moves to next record and return false if there is no more records
func (jl *KaminoJSONLoader) Next() bool {

	return len(jl.content) > jl.currentRow
}

//Load reads the next record and return it
func (jl *KaminoJSONLoader) Load() (types.Record, error) {
	if jl.currentRow > len(jl.content) {
		return nil, fmt.Errorf("no more data to read")
	}

	record := jl.content[jl.currentRow]
	jl.currentRow++
	return record, nil

}

//Close closes the datasource
func (jl *KaminoJSONLoader) Close() {
	//TODO: replace the following by the datasource.CloseFile()
	jl.ds.CloseFile()
}

//Name give the name of the destination
func (jl *KaminoJSONLoader) Name() string {
	return jl.name
}
