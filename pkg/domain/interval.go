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

// endpoint used in interval notation
type Inf struct{}
type MInf struct{}
type Int struct {
	v int
}

func (inf Inf) String() string   { return "+inf" }
func (minf MInf) String() string { return "-inf" }
func (i Int) String() string     { return fmt.Sprintf("%+d", i.v) }

// order of endpoint
func EPOrder(ep1, ep2 Interval) bool {
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

func EPMult(ep1, ep2 Interval) Interval {
	var ret Interval
	switch v1 := ep1.(type) {
	case MInf:
		switch v2 := ep2.(type) {
		case Inf:
			ret = MInf{}
		case MInf:
			ret = Inf{}
		case Int:
			if v2.v < 0 {
				ret = Inf{}
			} else if v2.v > 0 {
				ret = MInf{}
			} else {
				ret = Int{0}
			}
		}
	case Inf:
		switch v2 := ep2.(type) {
		case Inf:
			ret = Inf{}
		case MInf:
			ret = MInf{}
		case Int:
			if v2.v < 0 {
				ret = MInf{}
			} else if v2.v > 0 {
				ret = Inf{}
			} else {
				ret = Int{0}
			}
		}
	case Int:
		switch v2 := ep2.(type) {
		case Inf:
			if v1.v < 0 {
				ret = MInf{}
			} else if v1.v > 0 {
				ret = Inf{}
			} else {
				ret = Int{0}
			}
		case MInf:
			if v1.v < 0 {
				ret = Inf{}
			} else if v1.v > 0 {
				ret = MInf{}
			} else {
				ret = Int{0}
			}
		case Int:
			ret = Int{v1.v * v2.v}
		}
	}
	return ret
}

// interval - abstract value
type Bot struct{}
type Range struct {
	lbound Interval
	ubound Interval
}

// Abstract value to string
func (b Bot) String() string { return "Bot" }
func (r Range) String() string {
	return fmt.Sprintf("[%s, %s]", r.lbound.String(), r.ubound.String())
}

// abstraction from concrete value
func InterBot() Interval { return Bot{} }
func InterRange(lbound, ubound Interval) Interval {
	switch lbound.(type) {
	case Inf, MInf, Int:
		switch ubound.(type) {
		case Inf, MInf, Int:
			return Range{lbound, ubound}
		}
	}
	panic("InterRange: range define error (unreachable)")
}
func InterTop() Interval {
	return Range{MInf{}, Inf{}}
}

// order of interval
func InterOrder(i1, i2 Interval) bool {
	var ret bool
	switch i1.(type) {
	case Bot:
		ret = true
	case Range:
		switch i2.(type) {
		case Bot:
			ret = false
		case Range:
			i1_lbound, i1_ubound := i1.(Range).lbound, i1.(Range).ubound
			i2_lbound, i2_ubound := i2.(Range).lbound, i2.(Range).ubound
			ret = EPOrder(i2_lbound, i1_lbound) && EPOrder(i1_ubound, i2_ubound)
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
		i1_lbound, i1_ubound := i1.(Range).lbound, i1.(Range).ubound
		i2_lbound, i2_ubound := i2.(Range).lbound, i2.(Range).ubound
		var new_lbound, new_ubound Interval
		if EPOrder(i1_lbound, i2_lbound) {
			new_lbound = i1_lbound
		} else {
			new_lbound = i2_lbound
		}
		if EPOrder(i1_ubound, i2_ubound) {
			new_ubound = i2_ubound
		} else {
			new_ubound = i1_ubound
		}
		ret = InterRange(new_lbound, new_ubound)
	}
	return ret
}

func InterWiden(i1, i2 Interval) Interval {
	var ret Interval
	switch i1.(type) {
	case Bot:
		ret = i2
	case Range:
		switch i2.(type) {
		case Bot:
			ret = i1
		case Range:
			i1_lbound, i1_ubound := i1.(Range).lbound, i1.(Range).ubound
			i2_lbound, i2_ubound := i2.(Range).lbound, i2.(Range).ubound
			var new_lbound, new_ubound Interval
			if EPOrder(i1_lbound, i2_lbound) {
				new_lbound = i1_lbound
			} else {
				new_lbound = MInf{}
			}
			if EPOrder(i1_ubound, i2_ubound) {
				new_ubound = Inf{}
			} else {
				new_ubound = i1_ubound
			}
			ret = InterRange(new_lbound, new_ubound)
		}
	}
	return ret
}

func InterNarrow(i1, i2 Interval) Interval {
	var ret Interval
	switch i1.(type) {
	case Bot:
		ret = InterBot()
	case Range:
		switch i2.(type) {
		case Bot:
			ret = InterBot()
		case Range:
			i1_lbound, i1_ubound := i1.(Range).lbound, i1.(Range).ubound
			i2_lbound, i2_ubound := i2.(Range).lbound, i2.(Range).ubound
			var new_lbound, new_ubound Interval
			if i1_lbound == (MInf{}) {
				new_lbound = i2_lbound
			} else {
				new_lbound = i1_lbound
			}
			if i1_ubound == (Inf{}) {
				new_ubound = i2_ubound
			} else {
				new_ubound = i1_ubound
			}
			ret = InterRange(new_lbound, new_ubound)
		}
	}
	return ret
}

// binary operation
func InterPlus(i1, i2 Interval) Interval {
	var ret Interval
	switch i1.(type) {
	case Bot:
		ret = InterBot()
	case Range:
		switch i2.(type) {
		case Bot:
			ret = InterBot()
		case Range:
			i1_lbound, i1_ubound := i1.(Range).lbound, i1.(Range).ubound
			i2_lbound, i2_ubound := i2.(Range).lbound, i2.(Range).ubound
			var new_lbound, new_ubound Interval
			if i1_lbound == (MInf{}) {
				new_lbound = i1_lbound
			} else if i2_lbound == (MInf{}) {
				new_lbound = i2_lbound
			} else {
				new_lbound = Int{i1_lbound.(Int).v + i2_lbound.(Int).v}
			}
			if i1_ubound == (Inf{}) {
				new_ubound = i1_ubound
			} else if i2_ubound == (Inf{}) {
				new_ubound = i2_ubound
			} else {
				new_ubound = Int{i1_ubound.(Int).v + i2_ubound.(Int).v}
			}
			ret = InterRange(new_lbound, new_ubound)
		}
	}
	return ret
}

func InterMinus(i1, i2 Interval) Interval {
	var ret Interval
	switch i1.(type) {
	case Bot:
		ret = InterBot()
	case Range:
		switch i2.(type) {
		case Bot:
			ret = InterBot()
		case Range:
			i1_lbound, i1_ubound := i1.(Range).lbound, i1.(Range).ubound
			i2_lbound, i2_ubound := i2.(Range).lbound, i2.(Range).ubound
			var new_lbound, new_ubound Interval
			if i1_lbound == (MInf{}) {
				new_lbound = i1_lbound
			} else if i2_ubound == (Inf{}) {
				new_lbound = MInf{}
			} else {
				new_lbound = Int{i1_lbound.(Int).v - i2_ubound.(Int).v}
			}
			if i1_ubound == (Inf{}) {
				new_ubound = i1_ubound
			} else if i2_lbound == (MInf{}) {
				new_ubound = Inf{}
			} else {
				new_ubound = Int{i1_ubound.(Int).v - i2_lbound.(Int).v}
			}
			ret = InterRange(new_lbound, new_ubound)
		}
	}
	return ret
}

func InterMult(i1, i2 Interval) Interval {
	var ret Interval
	switch i1.(type) {
	case Bot:
		ret = InterBot()
	case Range:
		switch i2.(type) {
		case Bot:
			ret = InterBot()
		case Range:
			i1_lbound, i1_ubound := i1.(Range).lbound, i1.(Range).ubound
			i2_lbound, i2_ubound := i2.(Range).lbound, i2.(Range).ubound
			ep_list := [4]Interval{
				EPMult(i1_lbound, i2_lbound),
				EPMult(i1_lbound, i2_ubound),
				EPMult(i1_ubound, i2_lbound),
				EPMult(i1_ubound, i2_ubound),
			}

			new_lbound := ep_list[0]
			new_ubound := ep_list[0]
			for _, ep := range ep_list {
				if EPOrder(ep, new_lbound) {
					new_lbound = ep
				}
				if EPOrder(new_ubound, ep) {
					new_ubound = ep
				}
			}
			ret = InterRange(new_lbound, new_ubound)
		}
	}
	return ret
}
