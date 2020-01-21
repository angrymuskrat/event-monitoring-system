package data

import "fmt"

func (p *Point) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%.4f,%.4f\"", p.Lat, p.Lon)), nil
}
