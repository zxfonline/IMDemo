package atomic

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"sync/atomic"
)

//默认字节序 小端法
var (
	DefaultEndian = binary.LittleEndian
)

// Int32 is an atomic wrapper around an int32.
type Int32 struct{ v int32 }

// NewInt32 creates an Int32.
func NewInt32(i int32) *Int32 {
	return &Int32{i}
}

// Load atomically loads the wrapped value.
func (i *Int32) Load() int32 {
	return atomic.LoadInt32(&i.v)
}

// Add atomically adds to the wrapped int32 and returns the new value.
func (i *Int32) Add(n int32) int32 {
	return atomic.AddInt32(&i.v, n)
}

// Sub atomically subtracts from the wrapped int32 and returns the new value.
func (i *Int32) Sub(n int32) int32 {
	return atomic.AddInt32(&i.v, -n)
}

// Inc atomically increments the wrapped int32 and returns the new value.
func (i *Int32) Inc() int32 {
	return i.Add(1)
}

// Dec atomically decrements the wrapped int32 and returns the new value.
func (i *Int32) Dec() int32 {
	return i.Sub(1)
}

// CAS is an atomic compare-and-swap.
func (i *Int32) CAS(old, newz int32) bool {
	return atomic.CompareAndSwapInt32(&i.v, old, newz)
}

// Store atomically stores the passed value.
func (i *Int32) Store(n int32) {
	atomic.StoreInt32(&i.v, n)
}

// Swap atomically swaps the wrapped int32 and returns the old value.
func (i *Int32) Swap(n int32) int32 {
	return atomic.SwapInt32(&i.v, n)
}

func (i *Int32) MarshalJSON() ([]byte, error) {
	return i.MarshalBinary()
}
func (i *Int32) GobEncode() ([]byte, error) {
	return i.MarshalBinary()
}
func (i *Int32) MarshalText() ([]byte, error) {
	return i.MarshalBinary()
}
func (i *Int32) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	v := atomic.LoadInt32(&i.v)
	if err := binary.Write(&buf, DefaultEndian, v); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

func (i *Int32) UnmarshalJSON(data []byte) error {
	return i.UnmarshalBinary(data)
}
func (i *Int32) GobDecode(data []byte) error {
	return i.UnmarshalBinary(data)
}
func (i *Int32) UnmarshalText(data []byte) error {
	return i.UnmarshalBinary(data)
}
func (i *Int32) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return io.EOF
	}
	var x int32
	if err := binary.Read(bytes.NewBuffer(data), DefaultEndian, &x); err != nil {
		return err
	} else {
		atomic.StoreInt32(&i.v, x)
		return nil
	}
}

// Int64 is an atomic wrapper around an int64.
type Int64 struct{ v int64 }

// NewInt64 creates an Int64.
func NewInt64(i int64) *Int64 {
	return &Int64{i}
}

// Load atomically loads the wrapped value.
func (i *Int64) Load() int64 {
	return atomic.LoadInt64(&i.v)
}

// Add atomically adds to the wrapped int64 and returns the new value.
func (i *Int64) Add(n int64) int64 {
	return atomic.AddInt64(&i.v, n)
}

// Sub atomically subtracts from the wrapped int64 and returns the new value.
func (i *Int64) Sub(n int64) int64 {
	return atomic.AddInt64(&i.v, -n)
}

// Inc atomically increments the wrapped int64 and returns the new value.
func (i *Int64) Inc() int64 {
	return i.Add(1)
}

// Dec atomically decrements the wrapped int64 and returns the new value.
func (i *Int64) Dec() int64 {
	return i.Sub(1)
}

// CAS is an atomic compare-and-swap.
func (i *Int64) CAS(old, newz int64) bool {
	return atomic.CompareAndSwapInt64(&i.v, old, newz)
}

// Store atomically stores the passed value.
func (i *Int64) Store(n int64) {
	atomic.StoreInt64(&i.v, n)
}

// Swap atomically swaps the wrapped int64 and returns the old value.
func (i *Int64) Swap(n int64) int64 {
	return atomic.SwapInt64(&i.v, n)
}

func (i *Int64) MarshalJSON() ([]byte, error) {
	return i.MarshalBinary()
}
func (i *Int64) GobEncode() ([]byte, error) {
	return i.MarshalBinary()
}
func (i *Int64) MarshalText() ([]byte, error) {
	return i.MarshalBinary()
}
func (i *Int64) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	v := atomic.LoadInt64(&i.v)
	if err := binary.Write(&buf, DefaultEndian, v); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

func (i *Int64) UnmarshalJSON(data []byte) error {
	return i.UnmarshalBinary(data)
}
func (i *Int64) GobDecode(data []byte) error {
	return i.UnmarshalBinary(data)
}
func (i *Int64) UnmarshalText(data []byte) error {
	return i.UnmarshalBinary(data)
}
func (i *Int64) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return io.EOF
	}
	var x int64
	if err := binary.Read(bytes.NewBuffer(data), DefaultEndian, &x); err != nil {
		return err
	} else {
		atomic.StoreInt64(&i.v, x)
		return nil
	}
}

// Uint32 is an atomic wrapper around an uint32.
type Uint32 struct{ v uint32 }

// NewUint32 creates a Uint32.
func NewUint32(i uint32) *Uint32 {
	return &Uint32{i}
}

// Load atomically loads the wrapped value.
func (i *Uint32) Load() uint32 {
	return atomic.LoadUint32(&i.v)
}

// Add atomically adds to the wrapped uint32 and returns the new value.
func (i *Uint32) Add(n uint32) uint32 {
	return atomic.AddUint32(&i.v, n)
}

// Sub atomically subtracts from the wrapped uint32 and returns the new value.
func (i *Uint32) Sub(n uint32) uint32 {
	return atomic.AddUint32(&i.v, ^(n - 1))
}

// Inc atomically increments the wrapped uint32 and returns the new value.
func (i *Uint32) Inc() uint32 {
	return i.Add(1)
}

// Dec atomically decrements the wrapped int32 and returns the new value.
func (i *Uint32) Dec() uint32 {
	return i.Sub(1)
}

// CAS is an atomic compare-and-swap.
func (i *Uint32) CAS(old, newz uint32) bool {
	return atomic.CompareAndSwapUint32(&i.v, old, newz)
}

// Store atomically stores the passed value.
func (i *Uint32) Store(n uint32) {
	atomic.StoreUint32(&i.v, n)
}

// Swap atomically swaps the wrapped uint32 and returns the old value.
func (i *Uint32) Swap(n uint32) uint32 {
	return atomic.SwapUint32(&i.v, n)
}

func (i *Uint32) MarshalJSON() ([]byte, error) {
	return i.MarshalBinary()
}
func (i *Uint32) GobEncode() ([]byte, error) {
	return i.MarshalBinary()
}
func (i *Uint32) MarshalText() ([]byte, error) {
	return i.MarshalBinary()
}
func (i *Uint32) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	v := atomic.LoadUint32(&i.v)
	if err := binary.Write(&buf, DefaultEndian, v); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

func (i *Uint32) UnmarshalJSON(data []byte) error {
	return i.UnmarshalBinary(data)
}
func (i *Uint32) GobDecode(data []byte) error {
	return i.UnmarshalBinary(data)
}
func (i *Uint32) UnmarshalText(data []byte) error {
	return i.UnmarshalBinary(data)
}
func (i *Uint32) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return io.EOF
	}
	var x uint32
	if err := binary.Read(bytes.NewBuffer(data), DefaultEndian, &x); err != nil {
		return err
	} else {
		atomic.StoreUint32(&i.v, x)
		return nil
	}
}

// Uint64 is an atomic wrapper around a uint64.
type Uint64 struct{ v uint64 }

// NewUint64 creates a Uint64.
func NewUint64(i uint64) *Uint64 {
	return &Uint64{i}
}

// Load atomically loads the wrapped value.
func (i *Uint64) Load() uint64 {
	return atomic.LoadUint64(&i.v)
}

// Add atomically adds to the wrapped uint64 and returns the new value.
func (i *Uint64) Add(n uint64) uint64 {
	return atomic.AddUint64(&i.v, n)
}

// Sub atomically subtracts from the wrapped uint64 and returns the new value.
func (i *Uint64) Sub(n uint64) uint64 {
	return atomic.AddUint64(&i.v, ^(n - 1))
}

// Inc atomically increments the wrapped uint64 and returns the new value.
func (i *Uint64) Inc() uint64 {
	return i.Add(1)
}

// Dec atomically decrements the wrapped uint64 and returns the new value.
func (i *Uint64) Dec() uint64 {
	return i.Sub(1)
}

// CAS is an atomic compare-and-swap.
func (i *Uint64) CAS(old, newz uint64) bool {
	return atomic.CompareAndSwapUint64(&i.v, old, newz)
}

// Store atomically stores the passed value.
func (i *Uint64) Store(n uint64) {
	atomic.StoreUint64(&i.v, n)
}

// Swap atomically swaps the wrapped uint64 and returns the old value.
func (i *Uint64) Swap(n uint64) uint64 {
	return atomic.SwapUint64(&i.v, n)
}

func (i *Uint64) MarshalJSON() ([]byte, error) {
	return i.MarshalBinary()
}
func (i *Uint64) GobEncode() ([]byte, error) {
	return i.MarshalBinary()
}
func (i *Uint64) MarshalText() ([]byte, error) {
	return i.MarshalBinary()
}
func (i *Uint64) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	v := atomic.LoadUint64(&i.v)
	if err := binary.Write(&buf, DefaultEndian, v); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

func (i *Uint64) UnmarshalJSON(data []byte) error {
	return i.UnmarshalBinary(data)
}
func (i *Uint64) GobDecode(data []byte) error {
	return i.UnmarshalBinary(data)
}
func (i *Uint64) UnmarshalText(data []byte) error {
	return i.UnmarshalBinary(data)
}
func (i *Uint64) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return io.EOF
	}
	var x uint64
	if err := binary.Read(bytes.NewBuffer(data), DefaultEndian, &x); err != nil {
		return err
	} else {
		atomic.StoreUint64(&i.v, x)
		return nil
	}
}

// Bool is an atomic Boolean.
type Bool struct{ v uint32 }

// NewBool creates a Bool.
func NewBool(initial bool) *Bool {
	return &Bool{boolToInt(initial)}
}

// Load atomically loads the Boolean.
func (b *Bool) Load() bool {
	return truthy(atomic.LoadUint32(&b.v))
}

// Store atomically stores the passed value.
func (b *Bool) Store(newz bool) {
	atomic.StoreUint32(&b.v, boolToInt(newz))
}

// Swap sets the given value and returns the previous value.
func (b *Bool) Swap(newz bool) bool {
	return truthy(atomic.SwapUint32(&b.v, boolToInt(newz)))
}

// Toggle atomically negates the Boolean and returns the previous value.
func (b *Bool) Toggle() bool {
	return truthy(atomic.AddUint32(&b.v, 1) - 1)
}

func (b *Bool) MarshalJSON() ([]byte, error) {
	return b.MarshalBinary()
}
func (b *Bool) GobEncode() ([]byte, error) {
	return b.MarshalBinary()
}
func (b *Bool) MarshalText() ([]byte, error) {
	return b.MarshalBinary()
}
func (b *Bool) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	v := atomic.LoadUint32(&b.v)
	if err := binary.Write(&buf, DefaultEndian, v); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

func (b *Bool) UnmarshalJSON(data []byte) error {
	return b.UnmarshalBinary(data)
}
func (b *Bool) GobDecode(data []byte) error {
	return b.UnmarshalBinary(data)
}
func (b *Bool) UnmarshalText(data []byte) error {
	return b.UnmarshalBinary(data)
}
func (b *Bool) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return io.EOF
	}
	var x uint32
	if err := binary.Read(bytes.NewBuffer(data), DefaultEndian, &x); err != nil {
		return err
	} else {
		atomic.StoreUint32(&b.v, x)
		return nil
	}
}

func truthy(n uint32) bool {
	return n&1 == 1
}

func boolToInt(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}

// Float64 is an atomic wrapper around float64.
type Float64 struct {
	v uint64
}

// NewFloat64 creates a Float64.
func NewFloat64(f float64) *Float64 {
	return &Float64{math.Float64bits(f)}
}

// Load atomically loads the wrapped value.
func (f *Float64) Load() float64 {
	return math.Float64frombits(atomic.LoadUint64(&f.v))
}

// Store atomically stores the passed value.
func (f *Float64) Store(s float64) {
	atomic.StoreUint64(&f.v, math.Float64bits(s))
}

// CAS is an atomic compare-and-swap.
func (f *Float64) CAS(old, newz float64) bool {
	return atomic.CompareAndSwapUint64(&f.v, math.Float64bits(old), math.Float64bits(newz))
}

func (f *Float64) GobEncode() ([]byte, error) {
	return f.MarshalBinary()
}
func (f *Float64) MarshalJSON() ([]byte, error) {
	return f.MarshalBinary()
}
func (f *Float64) MarshalText() ([]byte, error) {
	return f.MarshalBinary()
}
func (f *Float64) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	v := atomic.LoadUint64(&f.v)
	if err := binary.Write(&buf, DefaultEndian, v); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

func (f *Float64) UnmarshalJSON(data []byte) error {
	return f.UnmarshalBinary(data)
}
func (f *Float64) GobDecode(data []byte) error {
	return f.UnmarshalBinary(data)
}
func (f *Float64) UnmarshalText(data []byte) error {
	return f.UnmarshalBinary(data)
}
func (f *Float64) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return io.EOF
	}
	var x uint64
	if err := binary.Read(bytes.NewBuffer(data), DefaultEndian, &x); err != nil {
		return err
	} else {
		atomic.StoreUint64(&f.v, x)
		return nil
	}
}

// Float32 is an atomic wrapper around float32.
type Float32 struct {
	v uint32
}

// NewFloat32 creates a Float32.
func NewFloat32(f float32) *Float32 {
	return &Float32{math.Float32bits(f)}
}

// Load atomically loads the wrapped value.
func (f *Float32) Load() float32 {
	return math.Float32frombits(atomic.LoadUint32(&f.v))
}

// Store atomically stores the passed value.
func (f *Float32) Store(s float32) {
	atomic.StoreUint32(&f.v, math.Float32bits(s))
}

// CAS is an atomic compare-and-swap.
func (f *Float32) CAS(old, newz float32) bool {
	return atomic.CompareAndSwapUint32(&f.v, math.Float32bits(old), math.Float32bits(newz))
}

func (f *Float32) MarshalJSON() ([]byte, error) {
	return f.MarshalBinary()
}
func (f *Float32) GobEncode() ([]byte, error) {
	return f.MarshalBinary()
}
func (f *Float32) MarshalText() ([]byte, error) {
	return f.MarshalBinary()
}
func (f *Float32) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	v := atomic.LoadUint32(&f.v)
	if err := binary.Write(&buf, DefaultEndian, v); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

func (f *Float32) GobDecode(data []byte) error {
	return f.UnmarshalBinary(data)
}
func (f *Float32) UnmarshalJSON(data []byte) error {
	return f.UnmarshalBinary(data)
}
func (f *Float32) UnmarshalText(data []byte) error {
	return f.UnmarshalBinary(data)
}
func (f *Float32) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return io.EOF
	}
	var x uint32
	if err := binary.Read(bytes.NewBuffer(data), DefaultEndian, &x); err != nil {
		return err
	} else {
		atomic.StoreUint32(&f.v, x)
		return nil
	}
}
