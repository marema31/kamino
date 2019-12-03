package datasource

// TmplValues structure use for template rendering to avoid exposing the datasource structure to the template
type TmplValues struct {
	Name        string
	Database    string
	User        string
	Password    string
	Schema      string
	Host        string
	Port        string
	Tags        []string
	Type        string
	Engine      string
	FilePath    string
	Transaction bool
	//TODO: add named tags
}

// FillTmplValues return a struct for template operation with value corresponding to the provided datasource
func (ds *Datasource) FillTmplValues() TmplValues {
	var tv TmplValues
	tv.Name = ds.name
	tv.Transaction = ds.transaction
	tv.Database = ds.database
	tv.User = ds.user
	tv.Password = ds.userPw
	tv.Schema = ds.schema
	tv.Host = ds.host
	tv.Port = ds.port
	tv.Tags = ds.tags
	tv.Type = TypeToString(ds.dstype)
	tv.Engine = EngineToString(ds.engine)
	tv.FilePath = ds.file.FilePath
	return tv
}
