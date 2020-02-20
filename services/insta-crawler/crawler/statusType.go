package crawler

import (
	"bytes"
	"encoding/json"
)

type StatusType int

const (
	InvalidStatus StatusType = iota
	RunningStatus
	FinishedStatus
	FailedStatus
	ToFix
)

func parseStatusType(s string) StatusType {
	switch s {
	case "running":
		return RunningStatus
	case "finished":
		return FinishedStatus
	case "failed":
		return FailedStatus
	default:
		return InvalidStatus
	}
}

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

func (s StatusType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(s.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (s *StatusType) UnmarshalJSON(b []byte) error {
	var in string
	err := json.Unmarshal(b, &in)
	if err != nil {
		return err
	}
	*s = parseStatusType(in)
	return nil
}
