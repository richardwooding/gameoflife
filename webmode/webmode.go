package webmode

import (
	"errors"
)

type WebMode uint8

const (
	Unknown WebMode = iota
	Live
	Static
)

func (m WebMode) String() string {
	switch m {
	case Live:
		return "live"
	case Static:
		return "static"
	default:
		return "unknown"
	}
}

var InvalidWebMode = errors.New("invalid web mode")

func (m WebMode) IsLive() bool {
	return m == Live
}

func (m WebMode) IsStatic() bool {
	return m == Static
}
func (m WebMode) IsUnknown() bool {
	return m == Unknown
}
func (m WebMode) IsValid() bool {
	return m == Live || m == Static
}

func ParseWebMode(s string) (WebMode, error) {
	switch s {
	case "live":
		return Live, nil
	case "static":
		return Static, nil
	default:
		return Unknown, InvalidWebMode
	}
}
