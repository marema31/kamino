package json

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/marema31/kamino/provider/common"
)

//KaminoJsonLoader specifc state for database Saver provider
type KaminoJSONLoader struct {
	file       io.ReadCloser
	content    []map[string]string
	currentRow int
}

//NewLoader open the encoding process on provider file, read the whole file, unMarshal it and return a Loader compatible object
func NewLoader(ctx context.Context, config map[string]string, file io.ReadCloser) (*KaminoJSONLoader, error) {
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	k := KaminoJSONLoader{file: file, content: nil, currentRow: 0}
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
func (jl *KaminoJSONLoader) Load() (common.Record, error) {
	if jl.currentRow > len(jl.content) {
		return nil, fmt.Errorf("no more data to read")
	}

	record := jl.content[jl.currentRow]
	jl.currentRow++
	return record, nil

}

//Close closes the datasource
func (jl *KaminoJSONLoader) Close() {
	jl.file.Close()
}
