package domain

import (
	"fmt"
)

/////////////////////
// Interval Domain //
/////////////////////
type Endpoint interface {
	EPString() string
}

type Inf struct{}
type MInf struct{}
type Int struct {
	v int
}

func (inf Inf) EPString() string   { return "+inf" }
func (minf MInf) EPString() string { return "-inf" }
func (i Int) EPString() string     { return fmt.Sprintf("%+d", i.v) }

func EPTop() Endpoint      { return Inf{} }
func EPBot() Endpoint      { return MInf{} }
func EPVal(v int) Endpoint { return Int{v} }

func EPOrder(ep1, ep2 Endpoint) bool {
	var ret bool
	switch v1 := ep1.(type) {
	case MInf:
		ret = true
	case Int:
		switch v2 := ep2.(type) {
		case MInf:
			ret = false
		case Int:
			return v1.v <= v2.v
		case Inf:
			ret = true
		}
	case Inf:
		switch ep2.(type) {
		case MInf, Int:
			ret = false
		case Inf:
			ret = true
		}
	}
	return ret
}

type Interval interface {
	InterString() string
}

type Bot struct{}
type Range struct {
	low  Endpoint
	high Endpoint
}

// Abstract value to string
func (b Bot) InterString() string { return "Bot" }
func (r Range) InterString() string {
	return fmt.Sprintf("[%s, %s]", r.low.EPString(), r.high.EPString())
}

// abstraction: int -> Interval
func InterBot() Interval { return Bot{} }
func InterRange(low, high Endpoint) Interval {
	return Range{low, high}
}
func InterTop() Interval {
	return Range{EPBot(), EPTop()}
}

// order
func InterOrder(i1, i2 Interval) bool {
	var ret bool
	switch v1 := i1.(type) {
	case Bot:
		ret = true
	case Range:
		switch v2 := i2.(type) {
		case Bot:
			ret = false
		case Range:
			i1_low, i1_high := v1.low, v1.high
			i2_low, i2_high := v2.low, v2.high
			ret = EPOrder(i2_low, i1_low) && EPOrder(i1_high, i2_high)
		}
	}
	return ret
}

// set operation
func InterJoin(i1, i2 Interval) Interval {
	var ret Interval
	if InterOrder(i1, i2) {
		ret = i2
	} else if InterOrder(i2, i1) {
		ret = i1
	} else {
		i1_low, i1_high := i1.(Range).low, i1.(Range).high
		i2_low, i2_high := i2.(Range).low, i2.(Range).high
		var new_low, new_high Endpoint
		if EPOrder(i1_low, i2_low) {
			new_low = i1_low
		} else {
			new_low = i2_low
		}
		if EPOrder(i1_high, i2_high) {
			new_high = i2_high
		} else {
			new_high = i1_high
		}
		ret = InterRange(new_low, new_high)
	}
	return ret
}

func InterWiden(i1, i2 Interval) Interval {
	var ret Interval
	switch v1 := i1.(type) {
	case Bot:
		ret = i2
	case Range:
		switch v2 := i2.(type) {
		case Bot:
			ret = i1
		case Range:
			i1_low, i1_high := v1.low, v1.high
			i2_low, i2_high := v2.low, v2.high
			var new_low, new_high Endpoint
			if EPOrder(i1_low, i2_low) {
				new_low = i1_low
			} else {
				new_low = EPBot()
			}
			if EPOrder(i1_high, i2_high) {
				new_high = EPTop()
			} else {
				new_high = i1_high
			}
			ret = InterRange(new_low, new_high)
		}
	}
	return ret
}

func InterNarrow(i1, i2 Interval) Interval {
	var ret Interval
	switch v1 := i1.(type) {
	case Bot:
		ret = InterBot()
	case Range:
		switch v2 := i2.(type) {
		case Bot:
			ret = InterBot()
		case Range:
			i1_low, i1_high := v1.low, v1.high
			i2_low, i2_high := v2.low, v2.high
			var new_low, new_high Endpoint
			if i1_low == EPBot() {
				new_low = i2_low
			} else {
				new_low = i1_low
			}
			if i1_high == EPTop() {
				new_high = i2_high
			} else {
				new_high = i1_high
			}
			ret = InterRange(new_low, new_high)
		}
	}
	return ret
}
