package scanner

// Bytes returns the input bytes that this scanner is conigured with.
func (s *Scanner) Bytes() []byte {
	return s.inp
}
