package pp

import (
	"fmt"
	"time"
)

type Appointment struct {
	session      *Session
	Timestamp    time.Time
	Practitioner Practitioner
	Status       Status
	blockIndex   string
}

func (a *Appointment) String() string {
	const layout = "on Jan 2 '06 at 3:04 pm"
	return fmt.Sprintf(
		"%s appointment with %s %s",
		a.Status,
		a.Practitioner,
		a.Timestamp.Format(layout))
}
