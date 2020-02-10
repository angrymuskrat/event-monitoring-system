package status

type HistoricBuilding struct {
	SessionID string
	Status    string
}

func (s HistoricBuilding) Get() Status {
	return s
}

func (s HistoricBuilding) String() string {
	return "historic building"
}
