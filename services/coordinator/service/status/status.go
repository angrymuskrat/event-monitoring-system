package status

type Status interface {
	Get() Status
	String() string
}
