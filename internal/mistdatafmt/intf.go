package mistdatafmt

// TODO; rewrite using mapstructure..

type MistDataFmtIntf interface {
	GetJsonKeyValue(key string) (interface{}, error)
	GetJsonKeyValueAsStr(key string) (string, error)
	GetJsonKeyValueAsFloat64(key string) (float64, error)
	GetJsonKeyValueAsInt64(key string) (int64, error)
}
