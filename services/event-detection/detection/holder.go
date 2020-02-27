package detection

type eventHolder struct {
	users map[string]bool
	posts map[string]bool
	tags  map[string]int
}
