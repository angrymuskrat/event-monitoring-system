package status

type Monitoring struct {
	SessionID        string
	CurrentTimestamp int64
	Status           string
}

func (s Monitoring) Get() Status {
	return s
}

func (s Monitoring) String() string {
	return "monitoring"
}
