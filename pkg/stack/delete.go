package stack

func (s *Stack) Delete() {
	s.mut.Lock()
	s.sta = nil
	s.mut.Unlock()
}
