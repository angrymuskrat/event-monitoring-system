package data

type Entity interface {
	GetID() string
	Marshal() ([]byte, error)
	Unmarshal(d []byte) error
}
