package schema

type DataType string

const (
	String    DataType = "string"
	Integer   DataType = "integer"
	Float     DataType = "float"
	Boolean   DataType = "boolean"
	Timestamp DataType = "timestamp"
)

type Column struct {
	Name     string
	Type     DataType
	Nullable bool
}

type Schema struct {
	Columns []Column
}