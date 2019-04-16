package provider

type Record map[string]string

type Saver interface {
	Save(Record) error
	//	Close()
}

type Loader interface {
	Next() bool
	Load() (Record, error)
	Close()
}
