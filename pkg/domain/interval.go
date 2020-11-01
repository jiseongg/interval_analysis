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
type Endpoint struct {
	epval int
}

func (inf Inf) String() string    { return "+inf" }
func (minf MInf) String() string  { return "-inf" }
func (i Endpoint) String() string { return fmt.Sprintf("%+d", i.epval) }

// order of endpoint
func EPSLE(ep1, ep2 Interval) bool {
	var ret bool
	switch v1 := ep1.(type) {
	case MInf:
		ret = true
	case Endpoint:
		switch v2 := ep2.(type) {
		case MInf:
			ret = false
		case Endpoint:
			return v1.epval <= v2.epval
		case Inf:
			ret = true
		}
	case Inf:
		switch ep2.(type) {
		case MInf, Endpoint:
			ret = false
		case Inf:
			ret = true
		}
	}
	return ret
}

func EPSLT(ep1, ep2 Interval) bool {
	var ret bool
	switch v1 := ep1.(type) {
	case MInf:
		ret = true
	case Endpoint:
		switch v2 := ep2.(type) {
		case MInf:
			ret = false
		case Endpoint:
			return v1.epval < v2.epval
		case Inf:
			ret = true
		}
	case Inf:
		switch ep2.(type) {
		case MInf, Endpoint:
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
		case Endpoint:
			if v2.epval < 0 {
				ret = Inf{}
			} else if v2.epval > 0 {
				ret = MInf{}
			} else {
				ret = Endpoint{0}
			}
		}
	case Inf:
		switch v2 := ep2.(type) {
		case Inf:
			ret = Inf{}
		case MInf:
			ret = MInf{}
		case Endpoint:
			if v2.epval < 0 {
				ret = MInf{}
			} else if v2.epval > 0 {
				ret = Inf{}
			} else {
				ret = Endpoint{0}
			}
		}
	case Endpoint:
		switch v2 := ep2.(type) {
		case Inf:
			if v1.epval < 0 {
				ret = MInf{}
			} else if v1.epval > 0 {
				ret = Inf{}
			} else {
				ret = Endpoint{0}
			}
		case MInf:
			if v1.epval < 0 {
				ret = Inf{}
			} else if v1.epval > 0 {
				ret = MInf{}
			} else {
				ret = Endpoint{0}
			}
		case Endpoint:
			ret = Endpoint{v1.epval * v2.epval}
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
	case Inf, MInf, Endpoint:
		switch ubound.(type) {
		case Inf, MInf, Endpoint:
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
			ret = EPSLE(i2_lbound, i1_lbound) && EPSLE(i1_ubound, i2_ubound)
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
		if EPSLE(i1_lbound, i2_lbound) {
			new_lbound = i1_lbound
		} else {
			new_lbound = i2_lbound
		}
		if EPSLE(i1_ubound, i2_ubound) {
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
			if EPSLT(i2_lbound, i1_lbound) {
				new_lbound = MInf{}
			} else {
				new_lbound = i1_lbound
			}
			if EPSLT(i1_ubound, i2_ubound) {
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
				new_lbound = Endpoint{i1_lbound.(Endpoint).epval + i2_lbound.(Endpoint).epval}
			}
			if i1_ubound == (Inf{}) {
				new_ubound = i1_ubound
			} else if i2_ubound == (Inf{}) {
				new_ubound = i2_ubound
			} else {
				new_ubound = Endpoint{i1_ubound.(Endpoint).epval + i2_ubound.(Endpoint).epval}
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
				new_lbound = Endpoint{i1_lbound.(Endpoint).epval - i2_ubound.(Endpoint).epval}
			}
			if i1_ubound == (Inf{}) {
				new_ubound = i1_ubound
			} else if i2_lbound == (MInf{}) {
				new_ubound = Inf{}
			} else {
				new_ubound = Endpoint{i1_ubound.(Endpoint).epval - i2_lbound.(Endpoint).epval}
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
				if EPSLE(ep, new_lbound) {
					new_lbound = ep
				}
				if EPSLE(new_ubound, ep) {
					new_ubound = ep
				}
			}
			ret = InterRange(new_lbound, new_ubound)
		}
	}
	return ret
}

func InterSLT(i1, i2 Interval) Interval {
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
			if EPSLT(i1_ubound, i2_lbound) {
				ret = InterRange(Endpoint{1}, Endpoint{1})
			} else if EPSLT(i2_ubound, i1_lbound) {
				ret = InterRange(Endpoint{0}, Endpoint{0})
			} else {
				ret = InterTop()
			}
		}
	}
	return ret
}
