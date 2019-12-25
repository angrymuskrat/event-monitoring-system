package data

type PushResponseType int

const (
	InvalidType PushResponseType = iota
	PostPushed
	DuplicatedPostId
	DBError
)

func CountPushRequestTypes() int {
	return 3
}

func ParseType(t int32) PushResponseType {
	switch t {
	case 0:
		return PostPushed
	case 1:
		return DuplicatedPostId
	case 2:
		return DBError
	default:
		return InvalidType
	}
}

func (t PushResponseType) Int32() int32 {
	switch t {
	case PostPushed:
		return 0
	case DuplicatedPostId:
		return 1
	case DBError:
		return 2
	default:
		return -1
	}
}