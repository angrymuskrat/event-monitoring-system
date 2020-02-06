package service

type StatusType int

const (
	InvalidStatus StatusType = iota
	RunningStatus
	FinishedStatus
	FailedStatus
)
