package safe

import (
	"sync"
	"sync/atomic"
	"testing"
)

// TestResourceBasic 验证构造器与值语义方法的基本行为。
func TestResourceBasic(t *testing.T) {
	// 按值构造
	r := NewResourceOf(10)
	if got := r.Get(); got != 10 {
		t.Fatalf("Get() = %v, want 10", got)
	}
	if got := r.Read(); got == nil || *got != 10 {
		t.Fatalf("Read() 指针内容 = %v, want 10", got)
	}
	if got := r.Ptr(); got == nil || *got != 10 {
		t.Fatalf("Ptr() 指针内容 = %v, want 10", got)
	}

	// Set / Reset
	r.Set(20)
	if got := r.Get(); got != 20 {
		t.Fatalf("after Set, Get() = %v, want 20", got)
	}
	r.Reset()
	if got := r.Get(); got != 0 {
		t.Fatalf("after Reset, Get() = %v, want 0", got)
	}
}

// TestResourceNewNil 验证 NewResource(nil) 存储 T 的零值且不 panic。
func TestResourceNewNil(t *testing.T) {
	r := NewResource[int](nil)
	if got := r.Get(); got != 0 {
		t.Fatalf("NewResource(nil).Get() = %v, want 0", got)
	}
}

// TestResourceGetNilSafe 验证内部为 nil 指针时 Get 返回零值不 panic。
func TestResourceGetNilSafe(t *testing.T) {
	r := &Resource[int]{}
	if got := r.Get(); got != 0 {
		t.Fatalf("empty Resource Get() = %v, want 0", got)
	}
	// Write(nil) 后 Get 仍返回零值
	r.Write(nil)
	if got := r.Get(); got != 0 {
		t.Fatalf("Write(nil).Get() = %v, want 0", got)
	}
}

// TestResourceSwap 验证 Swap/SwapVal 返回旧值。
func TestResourceSwap(t *testing.T) {
	r := NewResourceOf(1)

	old := r.SwapVal(2)
	if old != 1 {
		t.Fatalf("SwapVal 旧值 = %v, want 1", old)
	}
	if got := r.Get(); got != 2 {
		t.Fatalf("after SwapVal, Get() = %v, want 2", got)
	}

	newPtr := new(int)
	*newPtr = 3
	oldPtr := r.Swap(newPtr)
	if oldPtr == nil || *oldPtr != 2 {
		t.Fatalf("Swap 旧指针 = %v, want &2", oldPtr)
	}
	if got := r.Get(); got != 3 {
		t.Fatalf("after Swap, Get() = %v, want 3", got)
	}
}

// TestResourceCompareAndSwap 验证 CAS 成功/失败分支。
// 注意：atomic.Pointer.CompareAndSwap 比较的是指针身份而非值相等，
// 因此旧指针必须取自内部实际存储的指针。
func TestResourceCompareAndSwap(t *testing.T) {
	// 构造时保存内部存储的指针，CAS 用该指针才能成功。
	stored := new(int)
	*stored = 5
	r := NewResource(stored)

	replacement := 6
	if ok := r.CompareAndSwap(stored, &replacement); !ok {
		t.Fatalf("CAS 应当成功")
	}
	if got := r.Get(); got != 6 {
		t.Fatalf("after CAS, Get() = %v, want 6", got)
	}

	// 用旧指针（已被替换）再次 CAS 应当失败。
	stale := stored
	if ok := r.CompareAndSwap(stale, &replacement); ok {
		t.Fatalf("CAS 使用已失效旧指针应当失败")
	}

	// 用指向相同值但不同地址的指针也应失败（身份不等）。
	valueSame := 6
	if ok := r.CompareAndSwap(&valueSame, &replacement); ok {
		t.Fatalf("CAS 用不同地址但相同值的指针应当失败")
	}
}

// TestResourceUpdate 验证 Update 的原子修改。
func TestResourceUpdate(t *testing.T) {
	r := NewResourceOf(10)
	r.Update(func(v int) int { return v + 5 })
	if got := r.Get(); got != 15 {
		t.Fatalf("after Update, Get() = %v, want 15", got)
	}
}

// TestResourceUpdateConcurrent 并发验证 Update 的最终一致性。
func TestResourceUpdateConcurrent(t *testing.T) {
	r := NewResourceOf(0)
	const goroutines = 16
	const iters = 1000
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iters; j++ {
				r.Update(func(v int) int { return v + 1 })
			}
		}()
	}
	wg.Wait()
	if got := r.Get(); got != goroutines*iters {
		t.Fatalf("并发 Update 最终值 = %v, want %v", got, goroutines*iters)
	}
}

// TestResourceMixConcurrent 混合读写下验证快照操作与新值语义不 panic。
func TestResourceMixConcurrent(t *testing.T) {
	r := NewResourceOf(100)
	const goroutines = 8
	var wg sync.WaitGroup
	wg.Add(goroutines)
	var reads atomic.Int64
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 500; j++ {
				_ = r.Get()
				reads.Add(1)
				r.Update(func(v int) int { return v + 1 })
				_ = r.SwapVal(100)
			}
		}()
	}
	wg.Wait()
	if reads.Load() == 0 {
		t.Fatalf("未执行任何读操作")
	}
}

// ---------- 基准测试 ----------

// BenchmarkResourceRead 原子版读。
func BenchmarkResourceRead(b *testing.B) {
	r := NewResourceOf(1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = r.Get()
	}
}

// BenchmarkResourceWrite 原子版写。
func BenchmarkResourceWrite(b *testing.B) {
	r := NewResourceOf(1)
	val := 2
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r.Write(&val)
	}
}

// BenchmarkResourceUpdate 原子版 RMW 修改。
func BenchmarkResourceUpdate(b *testing.B) {
	r := NewResourceOf(0)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r.Update(func(v int) int { return v + 1 })
	}
}

// BenchmarkResourceRWMutexRead RWMutex 版读。
func BenchmarkResourceRWMutexRead(b *testing.B) {
	r := NewResourceRWMutex(1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = r.Read()
	}
}

// BenchmarkResourceRWMutexWrite RWMutex 版写。
func BenchmarkResourceRWMutexWrite(b *testing.B) {
	r := NewResourceRWMutex(1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r.Write(2)
	}
}

// BenchmarkResourceRWMutexUpdate RWMutex 版改。
func BenchmarkResourceRWMutexUpdate(b *testing.B) {
	r := NewResourceRWMutex(0)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r.Update(func(v int) int { return v + 1 })
	}
}