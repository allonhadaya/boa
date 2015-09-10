package pp

import (
	"fmt"
	"time"
)

type Appointment struct {
	session *Session
	// There's an interesting quirk here: the unix timestamp used on pocapoint is
	// always one hour earlier than the real-world appointment time. This is most
	// likely an implementation bug related to timezones.
	pptimestamp  int64
	practitioner Practitioner
	status       Status
	blockIndex   string
}

func (a *Appointment) RealTimestamp() time.Time {
	return time.Unix(a.pptimestamp, 0) //.Add(time.Hour)
}

func (a *Appointment) String() string {
	const layout = "on Jan 2 '06 at 3:04 pm"
	return fmt.Sprintf(
		"%s meeting with %s %s",
		a.status,
		a.practitioner,
		a.RealTimestamp().Format(layout))
}
