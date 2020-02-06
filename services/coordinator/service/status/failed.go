package status

type Failed struct {
	Error error
}

func (s Failed) Get() Status {
	return s
}

func (s Failed) String() string {
	return "failed"
}
