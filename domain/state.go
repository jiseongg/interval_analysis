package domain

//////////////////
// Memory State //
//////////////////
type State map[string]Interval

func EmptyState() State {
	return make(State)
}

func (s State) String() string {
	if len(s) == 0 {
		return "{ }"
	}
	var res string
	for k, v := range s {
		res = res + "\t" + k + " |-> " + v.String() + "\n"
	}
	return res
}

func (s *State) Bind(x string, v Interval) {
	(*s)[x] = v
}

func (s *State) Find(x string) Interval {
	v, ok := (*s)[x]
	if !ok {
		return InterBot()
	}
	return v
}
