package service

type StatusType int

const (
	InvalidStatus StatusType = iota
	RunningStatus
	FinishedStatus
	FailedStatus
)

func (s StatusType) String() string {
	switch s {
	case RunningStatus:
		return "running"
	case FinishedStatus:
		return "finished"
	case FailedStatus:
		return "failed"
	default:
		return ""
	}
}
