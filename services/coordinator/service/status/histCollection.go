package status

type HistoricCollection struct {
	SessionID      string
	PostsCollected int
	LocationsLeft  int
}

func (s HistoricCollection) Get() Status {
	return s
}

func (s HistoricCollection) String() string {
	return "historic collection"
}
