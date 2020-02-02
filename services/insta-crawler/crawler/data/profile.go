package data

import "encoding/json"

type Profile struct {
	ID             string
	Username       string
	FullName       string
	Biography      string
	FollowersCount int
	FollowsCount   int
	Verified       bool
	Private        bool
}

func (p *Profile) Marshal() ([]byte, error) {
	b, err := json.Marshal(p)
	return b, err
}

func (p *Profile) Unmarshal(d []byte) error {
	return json.Unmarshal(d, p)
}
