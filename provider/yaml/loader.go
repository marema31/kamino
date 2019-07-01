package yaml

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/marema31/kamino/config"
	"github.com/marema31/kamino/provider/common"
)

//KaminoYAMLLoader specifc state for database Saver provider
type KaminoYAMLLoader struct {
	file       io.ReadCloser
	name       string
	content    []map[string]string
	currentRow int
}

//NewLoader open the encoding process on provider file, read the whole file, unMarshal it and return a Loader compatible object
func NewLoader(ctx context.Context, loaderConfig config.SourceConfig, file io.ReadCloser) (*KaminoYAMLLoader, error) {
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	k := KaminoYAMLLoader{file: file, name: loaderConfig.File, content: nil, currentRow: 0}
	k.content = make([]map[string]string, 0)

	err = yaml.Unmarshal([]byte(byteValue), &k.content)
	if err != nil {
		return nil, err
	}

	return &k, nil
}

//Next moves to next record and return false if there is no more records
func (yl *KaminoYAMLLoader) Next() bool {

	return len(yl.content) > yl.currentRow
}

//Load reads the next record and return it
func (yl *KaminoYAMLLoader) Load() (common.Record, error) {
	if yl.currentRow > len(yl.content) {
		return nil, fmt.Errorf("no more data to read")
	}

	record := yl.content[yl.currentRow]
	yl.currentRow++
	return record, nil

}

//Close closes the datasource
func (yl *KaminoYAMLLoader) Close() {
	yl.file.Close()
}

//Name give the name of the destination
func (yl *KaminoYAMLLoader) Name() string {
	return yl.name
}
