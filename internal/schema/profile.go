package schema

type ColumnProfile struct {
	Name        string
	Type        DataType
	Count       int
	NullCount   int
	UniqueCount int

	Min    float64
	Max    float64
	Mean   float64
	StdDev float64

	Frequencies map[string]int

	IsID bool
}

type NumericStats struct {
	Count  int
	Min    float64
	Max    float64
	Mean   float64
	StdDev float64
}

type ConditionalProfile struct {
	ConditionColumn string
	ConditionValue  string
	Numeric         map[string]NumericStats
}

type DatasetProfile struct {
	RowCount int
	Columns  []ColumnProfile

	NumericData [][]float64

	ConditionalProfiles []ConditionalProfile
}
