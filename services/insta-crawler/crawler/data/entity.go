package data

type Entity interface {
	Marshal() ([]byte, error)
	Unmarshal(d []byte) error
}
