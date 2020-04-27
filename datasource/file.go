package datasource

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"

	"github.com/Masterminds/sprig/v3"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// load a dile type datasource from the viper configuration
func loadFileDatasource(recipePath string, filename string, v *viper.Viper, engine Engine, envVar map[string]string) (Datasource, error) {
	var ds Datasource
	ds.dstype = File
	ds.engine = engine
	ds.name = filename
	ds.file.Inline = v.GetString("inline")

	type tmplEnv struct {
		Environments map[string]string
	}

	data := tmplEnv{Environments: envVar}

	fileTmpl, err := template.New("file").Funcs(sprig.FuncMap()).Parse(v.GetString("file"))
	if err != nil {
		return Datasource{}, fmt.Errorf("parsing file provided")
	}

	var file bytes.Buffer
	if err = fileTmpl.Execute(&file, data); err != nil {
		return Datasource{}, fmt.Errorf("expanding file provided")
	}

	ds.file.FilePath = file.String()

	if ds.file.FilePath != "" && ds.file.FilePath != "-" && ds.file.FilePath[0] != '/' {
		ds.file.FilePath = filepath.Join(recipePath, ds.file.FilePath)
	}

	URLTmpl, err := template.New("URL").Funcs(sprig.FuncMap()).Parse(v.GetString("URL"))
	if err != nil {
		return Datasource{}, fmt.Errorf("parsing URL provided")
	}

	var URL bytes.Buffer
	if err = URLTmpl.Execute(&URL, data); err != nil {
		return Datasource{}, fmt.Errorf("expanding URL provided")
	}

	ds.file.URL = URL.String()

	if ds.file.FilePath == "" && ds.file.URL == "" && ds.file.Inline == "" {
		return Datasource{}, fmt.Errorf("no file path or URL provided")
	}

	ds.tags = v.GetStringSlice("tags")
	if len(ds.tags) == 0 {
		ds.tags = []string{""}
	}

	ds.file.Zip = v.GetBool("zip")
	ds.file.Gzip = v.GetBool("gzip")
	ds.file.ZippedExt = EngineToString(engine)

	return ds, nil
}

//OpenReadFile open and return a io.ReadCloser corresponding to the datasource to be used by providers
func (ds *Datasource) OpenReadFile(log *logrus.Entry) (io.ReadCloser, error) {
	logFile := log.WithField("engine", EngineToString(ds.engine))
	return ds.file.OpenReadFile(logFile)
}

//OpenWriteFile open and return a io.WriteCloser corresponding to the datasource to be used by providers
func (ds *Datasource) OpenWriteFile(log *logrus.Entry) (io.WriteCloser, error) {
	logFile := log.WithField("engine", EngineToString(ds.engine))
	return ds.file.OpenWriteFile(logFile)
}

//ResetFile close the file and remove the temporary file
func (ds *Datasource) ResetFile(log *logrus.Entry) error {
	logFile := log.WithField("engine", EngineToString(ds.engine))
	return ds.file.ResetFile(logFile)
}

//CloseFile close the file and rename the temporary file to real name (if exists)
func (ds *Datasource) CloseFile(log *logrus.Entry) error {
	logFile := log.WithField("engine", EngineToString(ds.engine))
	return ds.file.CloseFile(logFile)
}

// Stat returns os.FileInfo on the file of the datasource
func (ds *Datasource) Stat() (os.FileInfo, error) {
	return ds.file.Stat()
}
