package timer

import "time"

type Timer interface {
	Now() time.Time
}

type timer struct {
	loc *time.Location
}

func New(loc *time.Location) *timer {
	if loc == nil {
		loc = time.Local
	}
	return &timer{loc: loc}
}

func (t *timer) Now() time.Time {
	return time.Now().In(t.loc)
}
