package atomic

import (
	"sync/atomic"
)

// String is an atomic type-safe wrapper around atomic.Value for strings.
type String struct{ v atomic.Value }

// NewString creates a String.
func NewString(str string) *String {
	s := &String{}
	if str != "" {
		s.Store(str)
	}
	return s
}

// Load atomically loads the wrapped string.
func (s *String) Load() string {
	v := s.v.Load()
	if v == nil {
		return ""
	}
	return v.(string)
}

// Store atomically stores the passed string.
// Note: Converting the string to an interface{} to store in the atomic.Value
// requires an allocation.
func (s *String) Store(str string) {
	s.v.Store(str)
}

func (s *String) MarshalJSON() ([]byte, error) {
	return s.MarshalBinary()
}
func (s *String) GobEncode() ([]byte, error) {
	return s.MarshalBinary()
}
func (s *String) MarshalText() ([]byte, error) {
	return s.MarshalBinary()
}
func (s *String) MarshalBinary() ([]byte, error) {
	v := s.Load()
	return []byte(v), nil
}

func (s *String) UnmarshalJSON(data []byte) error {
	return s.UnmarshalBinary(data)
}
func (s *String) GobDecode(data []byte) error {
	return s.UnmarshalBinary(data)
}
func (s *String) UnmarshalText(data []byte) error {
	return s.UnmarshalBinary(data)
}
func (s *String) UnmarshalBinary(data []byte) error {
	s.Store(string(data))
	return nil
}
