package domain

import (
	"fmt"
)

/////////////////////
// Interval Domain //
/////////////////////
type Interval interface {
	String() string
}

type Bot struct{}
type Top struct{}
type Range struct {
	low  int
	high int
}
type LeftRange struct {
	high int
}
type RightRange struct {
	low int
}

// Abstract value to string
func (b Bot) String() string { return "Bot" }
func (t Top) String() string { return "[-inf, +inf]" }
func (r Range) String() string {
	return fmt.Sprintf("[%+d, %+d]", r.low, r.high)
}
func (lr LeftRange) String() string {
	return fmt.Sprintf("[-inf, %+d]", lr.high)
}
func (rr RightRange) String() string {
	return fmt.Sprintf("[%+d, +inf]", rr.low)
}

// abstraction: int -> Interval
func InterTop() Interval { return Top{} }
func InterBot() Interval { return Bot{} }
func InterRange(low, high int) Interval {
	return Range{low, high}
}
func InterLeftRange(high int) Interval {
	return LeftRange{high}
}
func InterRightRange(low int) Interval {
	return RightRange{low}
}

// order
func InterOrder(i1, i2 Interval) bool {
	var ret bool
	switch v1 := i1.(type) {
	case Bot:
		ret = true
	case Top:
		ret = i2 == Top{}
	case Range:
		switch v2 := i2.(type) {
		case Bot:
			ret = false
		case Top:
			ret = true
		case Range:
			ret = (v1.low >= v2.low && v1.high <= v2.high)
		case LeftRange:
			ret = (v1.high <= v2.high)
		case RightRange:
			ret = (v1.low >= v2.low)
		}
	case LeftRange:
		switch v2 := i2.(type) {
		case Bot, Range, RightRange:
			ret = false
		case Top:
			ret = true
		case LeftRange:
			ret = (v1.high <= v2.high)
		}
	case RightRange:
		switch v2 := i2.(type) {
		case Bot, Range, LeftRange:
			ret = false
		case Top:
			ret = true
		case RightRange:
			ret = (v1.low >= v2.low)
		}
	}

	return ret
}
