package ctype

import (
	"fmt"
	"strings"
	"time"
)

type Bool bool

func (t Bool) Unwrape() bool {
	return bool(t)
}

type Int int

func (t Int) Unwrape() int {
	return int(t)
}

type Float float64

func (t Float) Unwrape() float64 {
	return float64(t)
}

type String string

func (t String) Unwrape() string {
	return string(t)
}

type TimeDuration time.Duration

func (t TimeDuration) Unwrape() time.Duration {
	return time.Duration(t)
}

func (t TimeDuration) UnwrapeWithUnit(unit string) time.Duration {
	return time.Duration(t) * timeUnit(unit)
}

type Slice []string

func (t Slice) Unwrape() []string {
	return []string(t)
}

func (t Slice) UnwrapeAsString() (s string) {
	switch len(t) {
	case 0:
		s = ""
	case 1:
		s = t[0]
	default:
		s = t[0]
		for _, v := range t[1:] {
			s = fmt.Sprintf("%s, %s", s, v)
		}
	}
	return
}

func timeUnit(unit string) time.Duration {
	switch strings.ToLower(unit) {
	case "nanosecond", "ns", "nsec":
		return time.Nanosecond
	case "microsecond", "us", "usec":
		return time.Microsecond
	case "millisecond", "ms", "msec":
		return time.Millisecond
	case "second", "s", "sec":
		return time.Second
	case "minute", "m", "min":
		return time.Minute
	case "hour", "h", "hr":
		return time.Hour
	case "day", "d":
		return time.Hour * 24
	case "week", "wk", "w":
		return time.Hour * 7
	case "month", "mon":
		return time.Hour * 24 * 30
	case "year", "y", "yr":
		return time.Hour * 24 * 365
	default:
		return -1
	}
}
