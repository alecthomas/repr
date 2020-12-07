package internal

import (
	"time"
)

type StructA struct {
	String1 string
	string2 string
	Time1   time.Time
	time2   time.Time
}

func MakeStructA(s string, t time.Time) StructA {
	return StructA{
		String1: s,
		string2: s + "2",
		Time1:   t,
		time2:   t.Add(2 * time.Hour),
	}
}
