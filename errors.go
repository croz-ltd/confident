package confident

import "fmt"

type ConfidentFileNotReadableError struct {
	Path string
}

func (e ConfidentFileNotReadableError) Error() string {
	return fmt.Sprintf("Defined configuration file at location %s is not readable.", e.Path)
}

type ConfidentUnmarshallingError struct {
	Path           string
	UnmarshalError error
}

func (e ConfidentUnmarshallingError) Error() string {
	return fmt.Sprintf("Fail to unmarshall defined configuration file at location %s.", e.Path)
}

type ConfidentMarshallingError struct {
	Path         string
	MarshalError error
}

func (e ConfidentMarshallingError) Error() string {
	return fmt.Sprintf("Fail to marshall defined configuration file to location %s.", e.Path)
}

type ConfidentFileCreationError struct {
	Path          string
	CreationError error
}

func (e ConfidentFileCreationError) Error() string {
	return fmt.Sprintf("Fail to create new configuration file to location %s.", e.Path)
}

type ConfidentWriteError struct {
	Path       string
	WriteError error
}

func (e ConfidentWriteError) Error() string {
	return fmt.Sprintf("Fail to write data to configuration file to location %s.", e.Path)
}
