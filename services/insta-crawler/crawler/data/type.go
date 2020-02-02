package data

import (
	"bytes"
	"encoding/json"
)

type CrawlingType int

const (
	InvalidType CrawlingType = iota
	LocationsType
	ProfilesType
	InternalProfilesType
	StoriesType
)

func parseType(s string) CrawlingType {
	switch s {
	case "locations":
		return LocationsType
	case "profiles":
		return ProfilesType
	case "profiles-internal":
		return InternalProfilesType
	case "stories":
		return StoriesType
	default:
		return InvalidType
	}
}

func (s CrawlingType) String() string {
	switch s {
	case LocationsType:
		return "locations"
	case ProfilesType:
		return "profiles"
	case InternalProfilesType:
		return "profiles-internal"
	case StoriesType:
		return "stories"
	default:
		return ""
	}
}

func (s CrawlingType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(s.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (s *CrawlingType) UnmarshalJSON(b []byte) error {
	var in string
	err := json.Unmarshal(b, &in)
	if err != nil {
		return err
	}
	*s = parseType(in)
	return nil
}