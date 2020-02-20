package status

import "fmt"

type Failed struct {
	Error error
}

func (s Failed) Get() Status {
	return s
}

func (s Failed) String() string {
	return fmt.Sprintf("failed with error: %v", s.Error)
}
