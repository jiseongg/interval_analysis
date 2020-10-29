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

// point-wise order
func StateOrder(s1, s2 State) bool {
	for k, v1 := range s1 {
		v2, ok := s2[k]
		if !ok {
			v2 = InterBot()
		}
		if !InterOrder(v1, v2) {
			return false
		}
	}
	return true
}

// point-wise operation
func StateJoin(s1, s2 State) State {
	s3 := EmptyState()
	for k, v := range s2 {
		s3[k] = v
	}
	for k, v1 := range s1 {
		v2, ok := s3[k]
		if !ok {
			v2 = InterBot()
		}
		s3[k] = InterJoin(v1, v2)
	}
	return s3
}

func StateWiden(s1, s2 State) State {
	s3 := EmptyState()
	for k, v := range s2 {
		s3[k] = v
	}
	for k, v1 := range s1 {
		v2, ok := s3[k]
		if !ok {
			v2 = InterBot()
		}
		s3[k] = InterWiden(v1, v2)
	}
	return s3
}

func StateNarrow(s1, s2 State) State {
	s3 := EmptyState()
	for k, v := range s2 {
		s3[k] = v
	}
	for k, v1 := range s1 {
		v2, ok := s3[k]
		if !ok {
			v2 = InterBot()
		}
		s3[k] = InterNarrow(v1, v2)
	}
	return s3
}
