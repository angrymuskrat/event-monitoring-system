package status

type HistoricBuilding struct {
	GridsBuilt int
}

func (s HistoricBuilding) Get() Status {
	return s
}

func (s HistoricBuilding) String() string {
	return "historic building"
}
