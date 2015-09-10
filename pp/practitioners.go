package pp

import "fmt"

type Practitioner int

const (
	Rebecca   Practitioner = 74
	Elizabeth Practitioner = 75
	Rocio     Practitioner = 146
)

func (p Practitioner) String() string {
	switch p {
	case Rebecca:
		return "Rebecca"
	case Elizabeth:
		return "Elizabeth"
	case Rocio:
		return "Rocio"
	default:
		panic(fmt.Sprintf("There is no name defined for practitioner %d.", p))
	}
}
