package schema

type ColumnProfile struct {
	Name        string
	Type        DataType
	Count       int
	NullCount   int
	UniqueCount int
	Min         float64
	Max         float64
	Mean        float64
	Frequencies map[string]int
}

type DatasetProfile struct {
	RowCount int
	Columns  []ColumnProfile
}
