package utility

// UnmarshalFunc abstracts how the byte stream will be unmarshalled.
type UnmarshalFunc func([]byte, interface{}) error

// ComplexNumberForMarshal can be unmarshaled from byte streams.
type ComplexNumberForMarshal struct {
	Real		float64	`json:"real" yaml:"real"`
	Imaginary	float64	`json:"imaginary" yaml:"imaginary"`
}