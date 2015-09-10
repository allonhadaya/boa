package pp

import "fmt"

type Status int

const (
	Available Status = iota
	Booked
	Taken
)

func (s Status) String() string {
	switch s {
	case Available:
		return "available"
	case Booked:
		return "booked"
	case Taken:
		return "taken"
	default:
		panic(fmt.Sprintf("There is no name defined for status %d.", s))
	}
}
