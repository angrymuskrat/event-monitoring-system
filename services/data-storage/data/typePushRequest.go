package data

type PushRequestType int

const (
	InvalidType PushRequestType = iota
	PostPushed
	DuplicatedPostId
	DBError
)

func parseType(t int32) PushRequestType {
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

func (t PushRequestType) Int32() int32 {
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