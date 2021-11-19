package atomic

import (
	"runtime"
	"testing"
)

const _parallelism = 4
const _iterations = 1000

func TestStressInt32(t *testing.T) {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(_parallelism))
	atom := &Int32{0}
	for i := 0; i < _parallelism; i++ {
		go func() {
			for j := 0; j < _iterations; j++ {
				atom.Load()
				atom.Add(1)
				atom.Sub(2)
				atom.Inc()
				atom.Dec()
				atom.CAS(1, 0)
				atom.Swap(5)
				atom.Store(1)
				bb, _ := atom.MarshalBinary()
				atom.UnmarshalBinary(bb)
			}
		}()
	}
}

func TestStressInt64(t *testing.T) {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(_parallelism))
	atom := &Int64{0}
	for i := 0; i < _parallelism; i++ {
		go func() {
			for j := 0; j < _iterations; j++ {
				atom.Load()
				atom.Add(1)
				atom.Sub(2)
				atom.Inc()
				atom.Dec()
				atom.CAS(1, 0)
				atom.Swap(5)
				atom.Store(1)
				bb, _ := atom.MarshalBinary()
				atom.UnmarshalBinary(bb)
			}
		}()
	}
}

func TestStressUint32(t *testing.T) {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(_parallelism))
	atom := &Uint32{0}
	for i := 0; i < _parallelism; i++ {
		go func() {
			for j := 0; j < _iterations; j++ {
				atom.Load()
				atom.Add(1)
				atom.Sub(2)
				atom.Inc()
				atom.Dec()
				atom.CAS(1, 0)
				atom.Swap(5)
				atom.Store(1)
				bb, _ := atom.MarshalBinary()
				atom.UnmarshalBinary(bb)
			}
		}()
	}
}

func TestStressUint64(t *testing.T) {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(_parallelism))
	atom := &Uint64{0}
	for i := 0; i < _parallelism; i++ {
		go func() {
			for j := 0; j < _iterations; j++ {
				atom.Load()
				atom.Add(1)
				atom.Sub(2)
				atom.Inc()
				atom.Dec()
				atom.CAS(1, 0)
				atom.Swap(5)
				atom.Store(1)
				bb, _ := atom.MarshalBinary()
				atom.UnmarshalBinary(bb)
			}
		}()
	}
}

func TestStressBool(t *testing.T) {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(_parallelism))
	atom := NewBool(false)
	for i := 0; i < _parallelism; i++ {
		go func() {
			for j := 0; j < _iterations; j++ {
				atom.Load()
				atom.Store(false)
				atom.Swap(true)
				atom.Load()
				atom.Toggle()
				atom.Toggle()
				bb, _ := atom.MarshalBinary()
				atom.UnmarshalBinary(bb)
			}
		}()
	}
}

func TestStressString(t *testing.T) {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(_parallelism))
	atom := NewString("")
	for i := 0; i < _parallelism; i++ {
		go func() {
			for j := 0; j < _iterations; j++ {
				atom.Load()
				atom.Store("abc")
				atom.Load()
				atom.Store("def")
				bb, _ := atom.MarshalBinary()
				atom.UnmarshalBinary(bb)
			}
		}()
	}
}
