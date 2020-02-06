package status

type Monitoring struct {
	CurrentTimestamp string
}

func (s Monitoring) Get() Status {
	return s
}

func (s Monitoring) String() string {
	return "monitoring"
}
