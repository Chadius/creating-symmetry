package utility

// UnmarshalFunc abstracts how the byte stream will be unmarshalled.
type UnmarshalFunc func([]byte, interface{}) error
